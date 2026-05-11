// Throwaway STACKIT infra for testing stackit-nuke.
//
// Auth comes from STACKIT_SERVICE_ACCOUNT_KEY_PATH (same env var stackit-nuke
// uses). Run:
//
//   pulumi stack init dev
//   cp Pulumi.dev.yaml.example Pulumi.dev.yaml && $EDITOR Pulumi.dev.yaml
//   pulumi up
//
// Then exercise the nuker:
//   stackit-nuke run --config ../examples/compute-only.yaml --no-dry-run
//
// And clean any residual Pulumi state:
//   pulumi destroy

import * as pulumi from "@pulumi/pulumi";
import * as stackit from "@stackitcloud/pulumi-stackit";

// Explicit provider so we can opt into beta resources (getImageV2 is beta in
// v0.0.6). Same auth as the default — STACKIT_SERVICE_ACCOUNT_KEY_PATH env.
const provider = new stackit.Provider("stackit", {
    enableBetaResources: true,
});

const cfg = new pulumi.Config();
const projectId = cfg.require("projectId");
const region = cfg.get("region") ?? "eu01";
const availabilityZone = cfg.get("availabilityZone") ?? "eu01-1";
const machineType = cfg.get("machineType") ?? "g2i.1";
// Image IDs are tenant + region scoped, so resolve by name instead of
// hardcoding. Override by setting `imageId` in Pulumi config.
const imageIdOverride = cfg.get("imageId");
// g2i.* (Intel) needs distro=ubuntu. ARM machine types need distro=ubuntu-arm64.
const imageDistro = cfg.get("imageDistro") ?? "ubuntu";
const imageVersion = cfg.get("imageVersion") ?? "24.04";
const serverCount = cfg.getNumber("serverCount") ?? 2;

const resolvedImageId: pulumi.Output<string> = imageIdOverride
    ? pulumi.output(imageIdOverride)
    : stackit.getImageV2Output({
          projectId,
          filter: { distro: imageDistro, version: imageVersion },
      }, { provider }).apply(img => {
          if (!img.imageId) {
              throw new Error(
                  `no STACKIT image matched distro=${imageDistro} version=${imageVersion} in project ${projectId}`,
              );
          }
          return img.imageId;
      });

const labels = {
    "managed-by": "pulumi",
    "purpose": "stackit-nuke-dev",
};

// ---------------------------------------------------------------------------
// Network + interface
// ---------------------------------------------------------------------------
const network = new stackit.Network("dev-network", {
    projectId,
    name: "stackit-nuke-dev",
    routed: false,
    ipv4Prefix: "10.42.0.0/24",
    ipv4Gateway: "10.42.0.1",
    ipv4Nameservers: ["1.1.1.1", "8.8.8.8"],
    labels,
}, { provider });

// One NIC per server so each gets its own interface (servers can be attached
// to at most one interface at create time in this provider).
const nics = Array.from({ length: serverCount }).map((_, i) =>
    new stackit.NetworkInterface(`dev-nic-${i}`, {
        projectId,
        networkId: network.networkId,
        allowedAddresses: ["10.42.0.0/24"],
    }, { provider }),
);

// ---------------------------------------------------------------------------
// Servers
// ---------------------------------------------------------------------------
const servers = nics.map((nic, i) =>
    new stackit.Server(`dev-server-${i}`, {
        projectId,
        name: `stackit-nuke-dev-${i}`,
        bootVolume: {
            size: 16,
            sourceType: "image",
            sourceId: resolvedImageId,
        },
        availabilityZone,
        machineType,
        networkInterfaces: [nic.networkInterfaceId],
        labels,
    }, { provider }),
);

// ---------------------------------------------------------------------------
// Standalone volume (not attached) — gives ComputeVolume something to find
// when its lister is implemented.
// ---------------------------------------------------------------------------
const volume = new stackit.Volume("dev-volume", {
    projectId,
    name: "stackit-nuke-dev-extra",
    availabilityZone,
    size: 8,
    labels,
}, { provider });

// ---------------------------------------------------------------------------
// Outputs — handy for confirming what stackit-nuke should see.
// ---------------------------------------------------------------------------
export const networkId = network.networkId;
export const serverIds = pulumi.all(servers.map(s => s.serverId));
export const serverNames = pulumi.all(servers.map(s => s.name));
export const volumeId = volume.volumeId;
