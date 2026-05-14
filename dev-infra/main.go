// Throwaway STACKIT infra for testing stackit-nuke.
//
// Auth comes from STACKIT_SERVICE_ACCOUNT_KEY_PATH (same env var stackit-nuke
// uses). Run:
//
//	pulumi stack init dev
//	cp Pulumi.dev.yaml.example Pulumi.dev.yaml && $EDITOR Pulumi.dev.yaml
//	pulumi up
//
// Then exercise the nuker:
//
//	stackit-nuke run --config ../examples/compute-only.yaml --no-dry-run
//
// And clean any residual Pulumi state:
//
//	pulumi destroy
package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/stackitcloud/pulumi-stackit/sdk/go/stackit"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Explicit provider so we can opt into beta resources (GetImageV2 is
		// beta in v0.0.6). Same auth as the default — STACKIT_SERVICE_ACCOUNT_KEY_PATH env.
		provider, err := stackit.NewProvider(ctx, "stackit", &stackit.ProviderArgs{
			EnableBetaResources: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		cfg := config.New(ctx, "")
		projectId := cfg.Require("projectId")
		region := cfg.Get("region")
		if region == "" {
			region = "eu01"
		}
		availabilityZone := cfg.Get("availabilityZone")
		if availabilityZone == "" {
			availabilityZone = "eu01-1"
		}
		machineType := cfg.Get("machineType")
		if machineType == "" {
			machineType = "g2i.1"
		}
		// Image IDs are tenant + region scoped, so resolve by name instead of
		// hardcoding. Override by setting `imageId` in Pulumi config.
		imageIdOverride := cfg.Get("imageId")
		// g2i.* (Intel) needs distro=ubuntu. ARM machine types need distro=ubuntu-arm64.
		imageDistro := cfg.Get("imageDistro")
		if imageDistro == "" {
			imageDistro = "ubuntu"
		}
		imageVersion := cfg.Get("imageVersion")
		if imageVersion == "" {
			imageVersion = "24.04"
		}
		serverCount, err := cfg.TryInt("serverCount")
		if err != nil {
			serverCount = 2
		}
		_ = region

		var resolvedImageId pulumi.StringInput
		if imageIdOverride != "" {
			resolvedImageId = pulumi.String(imageIdOverride)
		} else {
			img := stackit.GetImageV2Output(ctx, stackit.GetImageV2OutputArgs{
				ProjectId: pulumi.String(projectId),
				Filter: &stackit.GetImageV2FilterArgs{
					Distro:  pulumi.String(imageDistro),
					Version: pulumi.String(imageVersion),
				},
			}, pulumi.Provider(provider))
			resolvedImageId = img.ImageId().ApplyT(func(id *string) (string, error) {
				if id == nil || *id == "" {
					return "", fmt.Errorf("no STACKIT image matched distro=%s version=%s in project %s", imageDistro, imageVersion, projectId)
				}
				return *id, nil
			}).(pulumi.StringOutput)
		}

		labels := pulumi.StringMap{
			"managed-by": pulumi.String("pulumi"),
			"purpose":    pulumi.String("stackit-nuke-dev"),
		}

		// -------------------------------------------------------------------
		// Network + interface
		// -------------------------------------------------------------------
		network, err := stackit.NewNetwork(ctx, "dev-network", &stackit.NetworkArgs{
			ProjectId:       pulumi.String(projectId),
			Name:            pulumi.String("stackit-nuke-dev"),
			Routed:          pulumi.Bool(false),
			Ipv4Prefix:      pulumi.String("10.42.0.0/24"),
			Ipv4Gateway:     pulumi.String("10.42.0.1"),
			Ipv4Nameservers: pulumi.StringArray{pulumi.String("1.1.1.1"), pulumi.String("8.8.8.8")},
			Labels:          labels,
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		// One NIC per server so each gets its own interface (servers can be
		// attached to at most one interface at create time in this provider).
		nics := make([]*stackit.NetworkInterface, serverCount)
		for i := 0; i < serverCount; i++ {
			nic, err := stackit.NewNetworkInterface(ctx, fmt.Sprintf("dev-nic-%d", i), &stackit.NetworkInterfaceArgs{
				ProjectId:        pulumi.String(projectId),
				NetworkId:        network.NetworkId,
				AllowedAddresses: pulumi.StringArray{pulumi.String("10.42.0.0/24")},
			}, pulumi.Provider(provider))
			if err != nil {
				return err
			}
			nics[i] = nic
		}

		// -------------------------------------------------------------------
		// Servers
		// -------------------------------------------------------------------
		servers := make([]*stackit.Server, serverCount)
		for i, nic := range nics {
			server, err := stackit.NewServer(ctx, fmt.Sprintf("dev-server-%d", i), &stackit.ServerArgs{
				ProjectId: pulumi.String(projectId),
				Name:      pulumi.String(fmt.Sprintf("stackit-nuke-dev-%d", i)),
				BootVolume: &stackit.ServerBootVolumeArgs{
					Size:       pulumi.Int(16),
					SourceType: pulumi.String("image"),
					SourceId:   resolvedImageId,
				},
				AvailabilityZone:  pulumi.String(availabilityZone),
				MachineType:       pulumi.String(machineType),
				NetworkInterfaces: pulumi.StringArray{nic.NetworkInterfaceId},
				Labels:            labels,
			}, pulumi.Provider(provider))
			if err != nil {
				return err
			}
			servers[i] = server
		}

		// -------------------------------------------------------------------
		// Standalone volume (not attached) — gives ComputeVolume something to
		// find when its lister is implemented.
		// -------------------------------------------------------------------
		volume, err := stackit.NewVolume(ctx, "dev-volume", &stackit.VolumeArgs{
			ProjectId:        pulumi.String(projectId),
			Name:             pulumi.String("stackit-nuke-dev-extra"),
			AvailabilityZone: pulumi.String(availabilityZone),
			Size:             pulumi.Int(8),
			Labels:           labels,
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		// -------------------------------------------------------------------
		// Outputs — handy for confirming what stackit-nuke should see.
		// -------------------------------------------------------------------
		serverIds := make(pulumi.StringArray, len(servers))
		serverNames := make(pulumi.StringArray, len(servers))
		for i, s := range servers {
			serverIds[i] = s.ServerId
			serverNames[i] = s.Name
		}
		ctx.Export("networkId", network.NetworkId)
		ctx.Export("serverIds", serverIds)
		ctx.Export("serverNames", serverNames)
		ctx.Export("volumeId", volume.VolumeId)
		return nil
	})
}
