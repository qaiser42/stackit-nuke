# Presets

Presets are named filter sets you can reference from multiple accounts.

```yaml
presets:
  keep-shared-infra:
    filters:
      Network:
        - property: Name
          value: shared-vpc
      DNSZone:
        - property: Name
          value: "*.corp.example.com"
          type: glob

accounts:
  "11111111-...":
    presets: [keep-shared-infra]
  "22222222-...":
    presets: [keep-shared-infra]
```

A preset is merged into the account's filter set at run time. You can combine multiple presets and add account-specific filters on top.
