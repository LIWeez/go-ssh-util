package gcp

import (
	"bufio"
	"fmt"
	"go-ssh-util/config"
	"go-ssh-util/ssh"
	"os"
	"os/exec"
	"strings"

	"github.com/trzsz/promptui"
)

type GCEInstance struct {
	Name        string
	Zone        string
	MachineType string
	InternalIP  string
	ExternalIP  string
	Status      string
}

type GCEInstanceConfig struct {
	Name        string
	Zone        string
	MachineType string
	ImageInfo   ImageInfo
	// Add other configuration fields as needed
}

type MachineTypeInfo struct {
	Name   string
	Zone   string
	CPU    string
	Memory string
}

func RunGetVMs() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := "gcloud compute instances list"
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func RunStartVM() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	selectedGCE, err := ChooseGCE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := fmt.Sprintf("gcloud compute instances start %s --zone=%s", selectedGCE.Name, selectedGCE.Zone)
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func RunStopVM() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	selectedGCE, err := ChooseGCE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	command := fmt.Sprintf("gcloud compute instances stop %s --zone=%s", selectedGCE.Name, selectedGCE.Zone)
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

}

func ChooseGCE() (GCEInstance, error) {
	cmd := exec.Command("gcloud", "compute", "instances", "list")

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return GCEInstance{}, fmt.Errorf("Error: %v\n", err)
	}

	// Convert the output to a string
	outputStr := string(output)
	// fmt.Println(outputStr)

	// Parse the output to extract instance details
	instances := parseGCEInstances(outputStr)

	// Create a prompt for selecting an instance
	prompt := promptui.Select{
		Label: "Select an instance:",
		Items: instances,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }} ({{ .Status }})",
			Active:   "\U0001F4BB {{ .Name | cyan }} ({{ .Status | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .Status | red }})",
			Selected: "{{ .Name | red | cyan }}",
			Details: `
	--------- Detail ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Type:" | faint }}	{{ .MachineType }}
	{{ "Zone:" | faint }}	{{ .Zone }}
	{{ "IP:" | faint }}	{{ .InternalIP }}
	{{ "Status:" | faint }}	{{ .Status }}`,
		},
		Size: 10,
	}

	// Show the prompt and get the selected instance
	index, _, err := prompt.Run()
	if err != nil {
		return GCEInstance{}, fmt.Errorf("Error: %v\n", err)
	}

	// Get the selected instance by index
	return instances[index], nil
}

// Function to parse the output of 'gcloud compute instances list'
func parseGCEInstances(output string) []GCEInstance {
	lines := strings.Split(output, "\n")
	var instances []GCEInstance

	// Skip the header line
	if len(lines) >= 1 {
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			// fmt.Println(fields)
			// fmt.Println(len(fields))
			if len(fields) >= 5 { // Make sure there are at least 6 fields
				instance := GCEInstance{
					Name:        fields[0],
					Zone:        fields[1],
					MachineType: fields[2],
					InternalIP:  fields[3],
				}

				// Check if there is an external IP field
				if len(fields) == 5 {
					instance.Status = fields[4]
				}

				// Check if there is a status field
				if len(fields) == 6 {
					instance.ExternalIP = fields[4]
					instance.Status = fields[5]
				}

				instances = append(instances, instance)
			}
		}
	}

	return instances
}

func RunCreateGCEInstance() {
	selectedHost, ExecutionMode, err := config.ChooseAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Retrieve the list of available zones
	regions, err := getAvailableRegions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a zone
	selectedRegion, err := chooseRegion(regions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	zones, err := getAvailableZones(selectedRegion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a zone
	selectedZone, err := chooseZone(zones)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Retrieve the list of available machine type groups
	// groups, err := getMachineTypeGroups(selectedZone)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	return
	// }
	selectedSeries, err := chooseMachineTypeSeries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a machine type group
	selectedGroup, err := chooseMachineTypeGroup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// List machine types for the selected group and zone
	machineTypes, err := listMachineTypes(selectedZone, selectedSeries, selectedGroup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Allow the user to choose a machine type
	selectedMachineType, err := chooseMachineType(machineTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	selectedImage, err := chooseImage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Prompt the user for GCE instance configuration
	config, err := promptForGCEInstanceConfig(selectedZone, selectedMachineType, selectedImage)
	fmt.Println(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Create the command using strings.Builder
	var cmdBuilder strings.Builder
	cmdBuilder.WriteString("gcloud compute instances create ")
	cmdBuilder.WriteString(config.Name)
	cmdBuilder.WriteString(" --zone=")
	cmdBuilder.WriteString(config.Zone)
	cmdBuilder.WriteString(" --machine-type=")
	cmdBuilder.WriteString(config.MachineType)
	cmdBuilder.WriteString(" --image-project=")
	cmdBuilder.WriteString(config.ImageInfo.Project)
	cmdBuilder.WriteString(" --image=")
	cmdBuilder.WriteString(config.ImageInfo.Name)

	// Get the final command string
	command := cmdBuilder.String()

	// Construct the 'gcloud' command to create the GCE instance
	// command := fmt.Sprintf("gcloud compute instances create %s --zone=%s --machine-type=%s --image-project=%s --image=%s", config.Name, config.Zone, config.MachineType, config.ImageInfo.Project, config.ImageInfo.Family)
	// Add other flags and parameters as needed

	// Execute the 'gcloud' command
	if ExecutionMode == 1 {
		ssh.ExecuteLocalCommand(command)
	} else {
		ssh.ExecuteRemoteCommand(command, fmt.Sprintf("%s@%s", selectedHost.User, selectedHost.Host), selectedHost.Port)
	}

	fmt.Println("GCE instance created successfully.")
}

func chooseMachineTypeSeries() (string, error) {
	series := []struct {
		Label       string
		Description string
	}{
		{"c3", "Intel Sapphire Rapids CPU"},
		{"e2", "根據可用性選擇CPU"},
		{"n2", "Intel Cascade Lake 和 Ice Lake CPU"},
		{"n2d", "AMD EPYC CPU"},
		{"t2a", "Ampere Altra ARM CPU"},
		{"t2d", "AMD EPYC Milan CPU"},
		{"n1", "Intel Skylake CPU"},
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F4BB {{ .Label | cyan }} ({{ .Description | red }})",
		Inactive: "  {{ .Label | cyan }} ({{ .Description | red }})",
		Selected: "\U0001F4BB {{ .Label | red | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a machine type series:",
		Items:     series,
		Templates: templates,
		Size:      10,
	}

	// Show the prompt and get the selected option
	selectedIndex, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Get the label of the selected option
	selectedLabel := series[selectedIndex].Label
	fmt.Println(selectedLabel)
	return selectedLabel, nil
}

// Function to allow the user to choose a machine type group
func chooseMachineTypeGroup() (string, error) {
	prompt := promptui.Select{
		Label: "Select a machine type group:",
		Items: []string{"standard", "cpu", "mem", "gpu"},
		Size:  10,
	}

	// Show the prompt and get the selected group
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	fmt.Println(result)
	return result, nil
}

// Function to list machine types for the selected group and zone
func listMachineTypes(zone, series, group string) ([]MachineTypeInfo, error) {
	// Construct the 'gcloud' command to list machine types with the specified zone filter
	cmd := exec.Command("gcloud", "compute", "machine-types", "list", "--filter=zone:"+zone)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return nil, err
	}
	var machineTypes []MachineTypeInfo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "n2-standard") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				machineType := MachineTypeInfo{
					Name:   fields[0],
					Zone:   fields[1],
					CPU:    fields[2],
					Memory: fields[3],
				}
				machineTypes = append(machineTypes, machineType)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command to finish:", err)
		return nil, err
	}

	if len(machineTypes) == 0 {
		fmt.Println("No 'n2-highcpu' machine types found in the specified zone.")
		return nil, err
	}

	return machineTypes, nil
}

// Function to allow the user to choose a machine type
func chooseMachineType(machineTypes []MachineTypeInfo) (MachineTypeInfo, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F622 {{ .Name | cyan }} (Zone: {{ .Zone }}, CPU: {{ .CPU }}, Memory: {{ .Memory }})",
		Inactive: "  {{ .Name | cyan }} (Zone: {{ .Zone }}, CPU: {{ .CPU }}, Memory: {{ .Memory }})",
		Selected: "\U0001F622 {{ .Name | red | cyan }} (Zone: {{ .Zone | red }}, CPU: {{ .CPU | red }}, Memory: {{ .Memory | red }})",
		Details: `
	--------- Detail ----------
	{{ "Name:" | faint }}	{{ .Name }}
	{{ "Zone:" | faint }}	{{ .Zone }}
	{{ "CPU:" | faint }}	{{ .CPU }}
	{{ "Memory:" | faint }}	{{ .Memory }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a Machine Type",
		Items:     machineTypes,
		Templates: templates,
		Size:      10,
	}

	// Show the prompt and get the selected machine type
	index, _, err := prompt.Run()
	if err != nil {
		return MachineTypeInfo{}, err
	}

	selectedMachineType := machineTypes[index]

	return selectedMachineType, nil
}

// Function to get the list of available zones
func getAvailableZones(region string) ([]string, error) {
	cmd := exec.Command("gcloud", "compute", "zones", "list", fmt.Sprintf("--filter=region:%s", region))

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Parse the output to extract zone names
	var zones []string
	lines := strings.Split(string(output), "\n")
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			zones = append(zones, strings.Fields(line)[0])
		}
	}
	fmt.Println(zones)
	return zones, nil
}

// Function to allow the user to choose a zone
func chooseZone(zones []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a zone:",
		Items: zones,
		Size:  10,
	}

	// Show the prompt and get the selected zone
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Function to get the list of available zones
func getAvailableRegions() ([]string, error) {
	cmd := exec.Command("gcloud", "compute", "regions", "list")

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Parse the output to extract zone names
	var regions []string
	lines := strings.Split(string(output), "\n")
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			regions = append(regions, strings.Fields(line)[0])
		}
	}
	fmt.Println(regions)
	return regions, nil
}

// Function to allow the user to choose a zone
func chooseRegion(regions []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a region:",
		Items: regions,
		Size:  10,
	}

	// Show the prompt and get the selected zone
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Function to prompt the user for GCE instance configuration
func promptForGCEInstanceConfig(selectedZone string, selectedMachineType MachineTypeInfo, selectedImage ImageInfo) (GCEInstanceConfig, error) {
	prompt := []*promptui.Prompt{
		{
			Label: "Enter instance name:",
		},
		// {
		// 	Label:   "Zone:",
		// 	Default: selectedZone, // Set the default zone to the selected one
		// },
		// {
		// 	Label:   "Choose machine type:",
		// 	Default: selectedMachineType,
		// },
		// Add prompts for other configuration fields as needed
	}

	config := GCEInstanceConfig{}

	for i, p := range prompt {
		result, err := p.Run()
		if err != nil {
			return GCEInstanceConfig{}, err
		}
		switch i {
		case 0:
			config.Name = result
			// Set other configuration fields based on prompts
		}

	}
	config.Zone = selectedZone
	config.MachineType = selectedMachineType.Name
	config.ImageInfo = selectedImage
	return config, nil
}

type ImageInfo struct {
	Name    string
	Project string
	Family  string
}

func chooseImage() (ImageInfo, error) {
	cmd := exec.Command("gcloud", "compute", "images", "list")

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running gcloud command:", err)
		os.Exit(1)
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Extract project and family information into a slice of ImageInfo structs
	imageInfoList := []ImageInfo{}
	for _, line := range lines[1:] { // Skip the header
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			name := fields[0]
			project := fields[1]
			family := fields[2]
			imageInfo := ImageInfo{Name: name, Project: project, Family: family}
			imageInfoList = append(imageInfoList, imageInfo)
		}
	}

	// Create a prompt to select a PROJECT
	projectPrompt := promptui.Select{
		Label: "Select a PROJECT",
		Items: getUniqueProjects(imageInfoList), // Get unique projects from the struct slice
	}

	_, selectedProject, err := projectPrompt.Run() // Execute the prompt to select a project

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// selectedProject := strings.TrimSpace(result)

	// Display associated FAMILY values for the selected PROJECT
	selectedFamilies := getFamiliesForProject(imageInfoList, selectedProject)

	var selectedImageInfo ImageInfo
	if len(selectedFamilies) > 0 {
		familyPrompt := promptui.Select{
			Label: "Select a FAMILY for " + selectedProject,
			Items: selectedFamilies, // Display families associated with the selected project
		}

		_, selectedFamily, err := familyPrompt.Run() // Let the user select a FAMILY

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Find the corresponding NAME for the selected FAMILY
		selectedName := getNameForFamily(imageInfoList, selectedProject, selectedFamily)

		// Print the selected PROJECT and FAMILY
		selectedImageInfo = ImageInfo{Name: selectedName, Project: selectedProject, Family: selectedFamily}
		// fmt.Println("You selected PROJECT:", selectedProject)
		// fmt.Println("You selected FAMILY:", selectedFamily)
		// fmt.Println("You selected Name:", selectedName)
		return selectedImageInfo, nil
	} else {
		fmt.Println("No FAMILY values found for", selectedProject)
		return selectedImageInfo, err
	}

}

// getUniqueProjects extracts unique PROJECT values from the struct slice
func getUniqueProjects(imageInfoList []ImageInfo) []string {
	projectsMap := make(map[string]bool)
	var projects []string
	for _, imageInfo := range imageInfoList {
		if _, exists := projectsMap[imageInfo.Project]; !exists {
			projectsMap[imageInfo.Project] = true
			projects = append(projects, imageInfo.Project)
		}
	}
	return projects
}

// getFamiliesForProject returns a slice of FAMILY values associated with a selected PROJECT
func getFamiliesForProject(imageInfoList []ImageInfo, selectedProject string) []string {
	var families []string
	for _, imageInfo := range imageInfoList {
		if imageInfo.Project == selectedProject {
			families = append(families, imageInfo.Family)
		}
	}
	return families
}

// getNameForFamily returns the NAME associated with a selected FAMILY for a given PROJECT
func getNameForFamily(imageInfoList []ImageInfo, selectedProject, selectedFamily string) string {
	for _, imageInfo := range imageInfoList {
		if imageInfo.Project == selectedProject && imageInfo.Family == selectedFamily {
			return imageInfo.Name
		}
	}
	return "No matching NAME found"
}
