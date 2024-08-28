package okta

import (
	"io"

	"github.com/okta/okta-cli-client/utils"
	"github.com/spf13/cobra"
)

var DeviceCmd = &cobra.Command{
	Use:   "device",
	Long:  "Manage DeviceAPI",
}

func NewDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
		Use:   "device",
		Long:  "Manage DeviceAPI",
	}
	return cmd
}

func init() {
    rootCmd.AddCommand(DeviceCmd)
}

var (
    
    
    
)

func NewListDevicesCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "lists",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.ListDevices(apiClient.GetConfig().Context)
            
            
            
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
	ListDevicesCmd := NewListDevicesCmd()
    DeviceCmd.AddCommand(ListDevicesCmd)
}

var (
    
    
            GetDevicedeviceId string
        
    
)

func NewGetDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "get",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.GetDevice(apiClient.GetConfig().Context, GetDevicedeviceId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&GetDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
    

	return cmd
}

func init() {
	GetDeviceCmd := NewGetDeviceCmd()
    DeviceCmd.AddCommand(GetDeviceCmd)
}

var (
    
    
            DeleteDevicedeviceId string
        
    
)

func NewDeleteDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "delete",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.DeleteDevice(apiClient.GetConfig().Context, DeleteDevicedeviceId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&DeleteDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
    

	return cmd
}

func init() {
	DeleteDeviceCmd := NewDeleteDeviceCmd()
    DeviceCmd.AddCommand(DeleteDeviceCmd)
}

var (
    
    
            ActivateDevicedeviceId string
        
            ActivateDevicedata string
        
    
)

func NewActivateDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "activate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.ActivateDevice(apiClient.GetConfig().Context, ActivateDevicedeviceId)
            
            
            if ActivateDevicedata != "" {
                req = req.Data(ActivateDevicedata)
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

    
    
        cmd.Flags().StringVarP(&ActivateDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
        cmd.Flags().StringVarP(&ActivateDevicedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	ActivateDeviceCmd := NewActivateDeviceCmd()
    DeviceCmd.AddCommand(ActivateDeviceCmd)
}

var (
    
    
            DeactivateDevicedeviceId string
        
            DeactivateDevicedata string
        
    
)

func NewDeactivateDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "deactivate",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.DeactivateDevice(apiClient.GetConfig().Context, DeactivateDevicedeviceId)
            
            
            if DeactivateDevicedata != "" {
                req = req.Data(DeactivateDevicedata)
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

    
    
        cmd.Flags().StringVarP(&DeactivateDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
        cmd.Flags().StringVarP(&DeactivateDevicedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	DeactivateDeviceCmd := NewDeactivateDeviceCmd()
    DeviceCmd.AddCommand(DeactivateDeviceCmd)
}

var (
    
    
            SuspendDevicedeviceId string
        
            SuspendDevicedata string
        
    
)

func NewSuspendDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "suspend",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.SuspendDevice(apiClient.GetConfig().Context, SuspendDevicedeviceId)
            
            
            if SuspendDevicedata != "" {
                req = req.Data(SuspendDevicedata)
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

    
    
        cmd.Flags().StringVarP(&SuspendDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
        cmd.Flags().StringVarP(&SuspendDevicedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	SuspendDeviceCmd := NewSuspendDeviceCmd()
    DeviceCmd.AddCommand(SuspendDeviceCmd)
}

var (
    
    
            UnsuspendDevicedeviceId string
        
            UnsuspendDevicedata string
        
    
)

func NewUnsuspendDeviceCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "unsuspend",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.UnsuspendDevice(apiClient.GetConfig().Context, UnsuspendDevicedeviceId)
            
            
            if UnsuspendDevicedata != "" {
                req = req.Data(UnsuspendDevicedata)
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

    
    
        cmd.Flags().StringVarP(&UnsuspendDevicedeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
        cmd.Flags().StringVarP(&UnsuspendDevicedata, "data", "", "", "")
        cmd.MarkFlagRequired("data")
        
    

	return cmd
}

func init() {
	UnsuspendDeviceCmd := NewUnsuspendDeviceCmd()
    DeviceCmd.AddCommand(UnsuspendDeviceCmd)
}

var (
    
    
            ListDeviceUsersdeviceId string
        
    
)

func NewListDeviceUsersCmd() *cobra.Command {
    cmd := &cobra.Command{
	    Use:   "listUsers",
	  
        RunE: func(cmd *cobra.Command, args []string) error {
            
            
            
            req := apiClient.DeviceAPI.ListDeviceUsers(apiClient.GetConfig().Context, ListDeviceUsersdeviceId)
            
            
            
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

    
    
        cmd.Flags().StringVarP(&ListDeviceUsersdeviceId, "deviceId", "", "", "")
        cmd.MarkFlagRequired("deviceId")
        
    

	return cmd
}

func init() {
	ListDeviceUsersCmd := NewListDeviceUsersCmd()
    DeviceCmd.AddCommand(ListDeviceUsersCmd)
}