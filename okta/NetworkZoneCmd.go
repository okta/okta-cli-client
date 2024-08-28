package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var NetworkZoneCmd = &cobra.Command{
	Use:   "networkZone",
	Long:  "Manage NetworkZoneAPI",
}

func NewNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "networkZone",
		Long:  "Manage NetworkZoneAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(NetworkZoneCmd)
}

var (
    
    
            CreateNetworkZonedata string
        
    
)

func NewCreateNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "create",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.CreateNetworkZone(apiClient.GetConfig().Context)
            
            
            if CreateNetworkZonedata != "" {
                req = req.Data(CreateNetworkZonedata)
            }
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&CreateNetworkZonedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	CreateNetworkZoneCmd := NewCreateNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(CreateNetworkZoneCmd)
}

var (
    
    
    
)

func NewListNetworkZonesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "lists",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.ListNetworkZones(apiClient.GetConfig().Context)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
    

	return cmd
}

func init() {
	ListNetworkZonesCmd := NewListNetworkZonesCmd()
    NetworkZoneCmd.AddCommand(ListNetworkZonesCmd)
}

var (
    
    
            GetNetworkZonezoneId string
        
    
)

func NewGetNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "get",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.GetNetworkZone(apiClient.GetConfig().Context, GetNetworkZonezoneId)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&GetNetworkZonezoneId, "zoneId", "", "", "")
        cmd.MarkFlagRequired("zoneId")
        
    

	return cmd
}

func init() {
	GetNetworkZoneCmd := NewGetNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(GetNetworkZoneCmd)
}

var (
    
    
            ReplaceNetworkZonezoneId string
        
            ReplaceNetworkZonedata string
        
    
)

func NewReplaceNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "replace",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.ReplaceNetworkZone(apiClient.GetConfig().Context, ReplaceNetworkZonezoneId)
            
            
            if ReplaceNetworkZonedata != "" {
                req = req.Data(ReplaceNetworkZonedata)
            }
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&ReplaceNetworkZonezoneId, "zoneId", "", "", "")
        cmd.MarkFlagRequired("zoneId")
        
        cmd.Flags().StringVarP(&ReplaceNetworkZonedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ReplaceNetworkZoneCmd := NewReplaceNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(ReplaceNetworkZoneCmd)
}

var (
    
    
            DeleteNetworkZonezoneId string
        
    
)

func NewDeleteNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "delete",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.DeleteNetworkZone(apiClient.GetConfig().Context, DeleteNetworkZonezoneId)
            
            
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&DeleteNetworkZonezoneId, "zoneId", "", "", "")
        cmd.MarkFlagRequired("zoneId")
        
    

	return cmd
}

func init() {
	DeleteNetworkZoneCmd := NewDeleteNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(DeleteNetworkZoneCmd)
}

var (
    
    
            ActivateNetworkZonezoneId string
        
            ActivateNetworkZonedata string
        
    
)

func NewActivateNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "activate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.ActivateNetworkZone(apiClient.GetConfig().Context, ActivateNetworkZonezoneId)
            
            
            if ActivateNetworkZonedata != "" {
                req = req.Data(ActivateNetworkZonedata)
            }
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&ActivateNetworkZonezoneId, "zoneId", "", "", "")
        cmd.MarkFlagRequired("zoneId")
        
        cmd.Flags().StringVarP(&ActivateNetworkZonedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ActivateNetworkZoneCmd := NewActivateNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(ActivateNetworkZoneCmd)
}

var (
    
    
            DeactivateNetworkZonezoneId string
        
            DeactivateNetworkZonedata string
        
    
)

func NewDeactivateNetworkZoneCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deactivate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.NetworkZoneAPI.DeactivateNetworkZone(apiClient.GetConfig().Context, DeactivateNetworkZonezoneId)
            
            
            if DeactivateNetworkZonedata != "" {
                req = req.Data(DeactivateNetworkZonedata)
            }
            
            resp, err := req.Execute()
            if err != nil {
                return err
            }
		    d, err := io.ReadAll(resp.Body)
		    if err != nil {
			    return err
		    }
		    utils.PrettyPrintByte(d)
            cmd.Println(string(d))
		    return nil
        },
    }

    
    
        cmd.Flags().StringVarP(&DeactivateNetworkZonezoneId, "zoneId", "", "", "")
        cmd.MarkFlagRequired("zoneId")
        
        cmd.Flags().StringVarP(&DeactivateNetworkZonedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	DeactivateNetworkZoneCmd := NewDeactivateNetworkZoneCmd()
    NetworkZoneCmd.AddCommand(DeactivateNetworkZoneCmd)
}