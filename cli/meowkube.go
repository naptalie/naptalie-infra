package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	nodeCount int
	namespace string
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "meowkube",
	Short: "Meowkube - A cat-themed K3s cluster manager",
	Long: `
 /\_/\  Meowkube: Your purr-fect K3s cluster manager
( o.o ) 
 > ^ <  Easily install, uninstall, and manage your K3s cluster
        with adorable cat-themed commands!`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Install command
var installCmd = &cobra.Command{
	Use:   "cuddle",
	Short: "Install a new K3s cluster",
	Long:  "Cuddle with a new K3s cluster by installing it on your machine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("? Meow! Starting to cuddle a new K3s cluster...")
		fmt.Printf("? Setting up a cluster with %d nodes...\n", nodeCount)
		
		// Execute K3s install script
		installCmd := exec.Command("bash", "-c", "curl -sfL https://get.k3s.io | sh -")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		
		err := installCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong while installing K3s: ", err)
			return
		}
		
		// If nodeCount > 1, provide instructions for adding more nodes
		if nodeCount > 1 {
			// Get K3s token
			tokenCmd := exec.Command("bash", "-c", "sudo cat /var/lib/rancher/k3s/server/node-token")
			token, err := tokenCmd.Output()
			if err != nil {
				fmt.Println("? Hiss! Couldn't retrieve the node token: ", err)
				return
			}
			
			// Get current IP
			ipCmd := exec.Command("hostname", "-I")
			ip, err := ipCmd.Output()
			if err != nil {
				fmt.Println("? Hiss! Couldn't determine IP address: ", err)
				return
			}
			
			fmt.Println("? Purr! Master node is ready!")
			fmt.Println("? To add more nodes, run this command on each node:")
			fmt.Printf("curl -sfL https://get.k3s.io | K3S_URL=https://%s:6443 K3S_TOKEN=%s sh -\n", 
				strings.TrimSpace(string(ip)), strings.TrimSpace(string(token)))
		}
		
		fmt.Println("? Meow! Your K3s cluster is ready! Use 'meowkube pounce' to check its status!")
	},
}

// Uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "scratch",
	Short: "Uninstall the K3s cluster",
	Long:  "Scratch the K3s cluster away from your system",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("? Meow! Starting to scratch away the K3s cluster...")
		
		// Execute K3s uninstall script
		uninstallCmd := exec.Command("bash", "-c", "/usr/local/bin/k3s-uninstall.sh")
		uninstallCmd.Stdout = os.Stdout
		uninstallCmd.Stderr = os.Stderr
		
		err := uninstallCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong while uninstalling K3s: ", err)
			return
		}
		
		fmt.Println("? Purr! Your K3s cluster has been scratched away successfully!")
	},
}

// Get command (kubectl get wrapper)
var getCmd = &cobra.Command{
	Use:   "peek [resource]",
	Short: "Get resources (kubectl get wrapper)",
	Long:  "Peek at your cluster resources (wraps kubectl get)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		fmt.Printf("? Meow! Peeking at %s in namespace %s...\n", resource, namespace)
		
		kubectlArgs := []string{"get", resource}
		if namespace != "" && namespace != "all" {
			kubectlArgs = append(kubectlArgs, "-n", namespace)
		} else if namespace == "all" {
			kubectlArgs = append(kubectlArgs, "--all-namespaces")
		}
		
		// Add any additional arguments passed
		if len(args) > 1 {
			kubectlArgs = append(kubectlArgs, args[1:]...)
		}
		
		// Execute kubectl command
		kubectlCmd := exec.Command("kubectl", kubectlArgs...)
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr
		
		err := kubectlCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong: ", err)
			return
		}
	},
}

// Describe command (kubectl describe wrapper)
var describeCmd = &cobra.Command{
	Use:   "inspect [resource] [name]",
	Short: "Describe resources (kubectl describe wrapper)",
	Long:  "Inspect your cluster resources in detail (wraps kubectl describe)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		name := args[1]
		fmt.Printf("? Meow! Inspecting %s/%s in namespace %s...\n", resource, name, namespace)
		
		kubectlArgs := []string{"describe", resource, name}
		if namespace != "" && namespace != "all" {
			kubectlArgs = append(kubectlArgs, "-n", namespace)
		} else if namespace == "all" {
			kubectlArgs = append(kubectlArgs, "--all-namespaces")
		}
		
		// Execute kubectl command
		kubectlCmd := exec.Command("kubectl", kubectlArgs...)
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr
		
		err := kubectlCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong: ", err)
			return
		}
	},
}

// Delete command (kubectl delete wrapper)
var deleteCmd = &cobra.Command{
	Use:   "swat [resource] [name]",
	Short: "Delete resources (kubectl delete wrapper)",
	Long:  "Swat away your cluster resources (wraps kubectl delete)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		name := args[1]
		fmt.Printf("? Meow! Swatting at %s/%s in namespace %s...\n", resource, name, namespace)
		
		kubectlArgs := []string{"delete", resource, name}
		if namespace != "" && namespace != "all" {
			kubectlArgs = append(kubectlArgs, "-n", namespace)
		} else if namespace == "all" {
			kubectlArgs = append(kubectlArgs, "--all-namespaces")
		}
		
		// Execute kubectl command
		kubectlCmd := exec.Command("kubectl", kubectlArgs...)
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr
		
		err := kubectlCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong: ", err)
			return
		}
		
		fmt.Printf("? Purr! Successfully swatted away %s/%s!\n", resource, name)
	},
}

// Apply command (kubectl apply wrapper)
var applyCmd = &cobra.Command{
	Use:   "pounce -f [filename]",
	Short: "Apply resources (kubectl apply wrapper)",
	Long:  "Pounce on your manifest files to apply them to the cluster (wraps kubectl apply)",
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		if filename == "" {
			fmt.Println("? Hiss! You must specify a filename with -f")
			return
		}
		
		fmt.Printf("? Meow! Pouncing on %s in namespace %s...\n", filename, namespace)
		
		kubectlArgs := []string{"apply", "-f", filename}
		if namespace != "" && namespace != "all" {
			kubectlArgs = append(kubectlArgs, "-n", namespace)
		}
		
		// Execute kubectl command
		kubectlCmd := exec.Command("kubectl", kubectlArgs...)
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr
		
		err := kubectlCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong: ", err)
			return
		}
		
		fmt.Printf("? Purr! Successfully pounced on %s!\n", filename)
	},
}

// Logs command (kubectl logs wrapper)
var logsCmd = &cobra.Command{
	Use:   "meow [pod]",
	Short: "View pod logs (kubectl logs wrapper)",
	Long:  "Let your pods meow their logs to you (wraps kubectl logs)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pod := args[0]
		follow, _ := cmd.Flags().GetBool("follow")
		
		fmt.Printf("? Meow! Listening to %s in namespace %s...\n", pod, namespace)
		
		kubectlArgs := []string{"logs", pod}
		if namespace != "" {
			kubectlArgs = append(kubectlArgs, "-n", namespace)
		}
		
		if follow {
			kubectlArgs = append(kubectlArgs, "-f")
			fmt.Println("? Keeping an ear out for new meows... (Ctrl+C to stop)")
		}
		
		// Execute kubectl command
		kubectlCmd := exec.Command("kubectl", kubectlArgs...)
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr
		
		err := kubectlCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something went wrong: ", err)
			return
		}
	},
}

// Cluster status command
var statusCmd = &cobra.Command{
	Use:   "purr",
	Short: "Check the status of your K3s cluster",
	Long:  "Make your K3s cluster purr and tell you its status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("? Meow! Checking your K3s cluster's purr...")
		
		// Check node status
		nodeCmd := exec.Command("kubectl", "get", "nodes")
		nodeCmd.Stdout = os.Stdout
		nodeCmd.Stderr = os.Stderr
		
		err := nodeCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Something's wrong with your cluster: ", err)
			return
		}
		
		// Check pod status
		fmt.Println("\n? Checking on system kittens (pods)...")
		podCmd := exec.Command("kubectl", "get", "pods", "--all-namespaces")
		podCmd.Stdout = os.Stdout
		podCmd.Stderr = os.Stderr
		
		err = podCmd.Run()
		if err != nil {
			fmt.Println("? Hiss! Couldn't get pod information: ", err)
			return
		}
		
		fmt.Println("\n? Purr! Your K3s cluster is running smoothly!")
	},
}
// TUI command
var tuiCmd = &cobra.Command{
	Use:   "purr-tui",
	Short: "Start the interactive MeowTUI",
	Long:  "Start MeowTUI - An interactive pink cat-themed terminal UI for managing your K3s cluster",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := SimpleMeowTUI()
		if err != nil {
			fmt.Println("? Hiss! Error running TUI:", err)
		}
	},
}
func init() {
	// Add commands to root command
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(tuiCmd)	
	// Flags for install command
	installCmd.Flags().IntVarP(&nodeCount, "nodes", "n", 1, "Number of nodes in the cluster")
	
	// Flags for commands that use namespace
	getCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (use 'all' for --all-namespaces)")
	describeCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (use 'all' for --all-namespaces)")
	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (use 'all' for --all-namespaces)")
	applyCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")
	logsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace")
	
	// Apply command requires a filename
	applyCmd.Flags().StringP("filename", "f", "", "Filename to apply")
	applyCmd.MarkFlagRequired("filename")
	
	// Logs command can follow logs
	logsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("? Catastrophe! ", err)
		os.Exit(1)
	}
}
