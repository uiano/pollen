package gopher

import (
    "fmt"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/spf13/viper"
)

func GetClient() *gophercloud.ServiceClient {
    opts := gophercloud.AuthOptions{
        IdentityEndpoint: viper.GetString("IKT_STACK_IDENTITY_ENDPOINT"),
        Username:         viper.GetString("IKT_STACK_USERNAME"),
        Password:         viper.GetString("IKT_STACK_PASSWORD"),
        DomainName:       viper.GetString("IKT_STACK_DOMAIN_NAME"),
        TenantName:       viper.GetString("IKT_STACK_TENANT_NAME"),
    }

    provider, err := openstack.AuthenticatedClient(opts)
    if err != nil {
        fmt.Println("Provider err: ", err)
    }

    opts1 := gophercloud.EndpointOpts{Region: viper.GetString("IKT_STACK_REGION")}

    client, err := openstack.NewComputeV2(provider, opts1)
    if err != nil {
        fmt.Println("Client err: ", err)
    }

    return client
}