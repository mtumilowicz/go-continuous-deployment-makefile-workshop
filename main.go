package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

// -commit f7849357b169da3d0f446b2717a3eef644159fdd -chart-dir ./helm -image-name helm-workshop -namespace default -release-name helmworkshopchart
// -commit 292e275cacf0238ec0e3d76e8c4948a02c051fc7 -chart-dir ./helm -image-name helm-workshop -namespace default -release-name helmworkshopchart
func main() {
	// Parse command-line arguments
	commitHash := flag.String("commit", "", "Git commit hash (required)")
	chartDir := flag.String("chart-dir", "./helm", "Path to Helm chart directory")
	imageName := flag.String("image-name", "", "Name of the Docker image (required)")
	namespace := flag.String("namespace", "default", "Kubernetes namespace")
	releaseName := flag.String("release-name", "", "Name of the Helm release (required)")
	flag.Parse()

	// Validate required flags
	if *commitHash == "" {
		fmt.Println("Commit hash is required.")
		flag.Usage()
		os.Exit(1)
	}
	if *imageName == "" {
		fmt.Println("Image name is required.")
		flag.Usage()
		os.Exit(1)
	}
	if *releaseName == "" {
		fmt.Println("Release name is required.")
		flag.Usage()
		os.Exit(1)
	}

	// Set repository URL and clone directory
	repoURL := "https://github.com/mtumilowicz/helm-workshop"
	cloneDir := "./helm-workshop"

	// Clone the repository
	fmt.Printf("Cloning repository from %s...\n", repoURL)
	err := cloneRepository(repoURL, cloneDir)
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking out commit: %s\n", *commitHash)
	err = checkoutCommit(*commitHash, cloneDir)
	if err != nil {
		fmt.Printf("Error checking out commit: %v\n", err)
		os.Exit(1)
	}

	// Clean build artifacts
	fmt.Printf("Cleaning build artifacts...\n")
	err = cleanBuild(cloneDir)
	if err != nil {
		fmt.Printf("Error cleaning build artifacts: %v\n", err)
		os.Exit(1)
	}

	// Build Docker image using Gradle
	fmt.Printf("Building Docker image with tag: %s\n", *commitHash)
	err = buildDockerImage(*commitHash, *imageName, cloneDir)
	if err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		os.Exit(1)
	}

	// Upgrade Helm chart
	fmt.Printf("Upgrading Helm chart with image tag: %s\n", *commitHash)
	err = upgradeHelmChart(*commitHash, *chartDir, *namespace, *releaseName, cloneDir)
	if err != nil {
		fmt.Printf("Error upgrading Helm chart: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Deployment successful.")
}

// cloneRepository clones the Git repository to the specified directory.
func cloneRepository(repoURL, cloneDir string) error {
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// checkoutCommit checks out the specified commit hash in the working directory.
func checkoutCommit(commitHash, cloneDir string) error {
	cmd := exec.Command("git", "checkout", commitHash)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// cleanBuild removes build artifacts from the working directory.
func cleanBuild(cloneDir string) error {
	cmd := exec.Command("./gradlew", "clean")
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// buildDockerImage builds the Docker image using Gradle with the given commit hash as the tag.
func buildDockerImage(commitHash, imageName, cloneDir string) error {
	// Format the image tag
	tag := fmt.Sprintf("%s:%s", imageName, commitHash)
	// Run Gradle with the image name and tag
	cmd := exec.Command("./gradlew", "bootBuildImage", "--imageName="+tag)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// upgradeHelmChart upgrades the Helm chart with the given image tag, chart directory, and release name.
func upgradeHelmChart(commitHash, chartDir, namespace, releaseName, cloneDir string) error {
	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartDir, "--set", fmt.Sprintf("deployment.image.version=%s", commitHash), "--namespace", namespace)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
