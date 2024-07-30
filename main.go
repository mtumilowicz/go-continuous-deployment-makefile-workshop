package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Parse command-line arguments
	//commitHash := flag.String("commit", "292e275cacf0238ec0e3d76e8c4948a02c051fc7", "Git commit hash")
	commitHash := flag.String("commit", "f7849357b169da3d0f446b2717a3eef644159fdd", "Git commit hash")
	chartDir := flag.String("chart-dir", "./helm", "Path to Helm chart directory")
	imageName := flag.String("image-name", "helm-workshop", "Name of the Docker image")
	namespace := flag.String("namespace", "default", "Kubernetes namespace")
	releaseName := flag.String("release-name", "helmworkshopchart", "Name of the Helm release")
	flag.Parse()

	if *commitHash == "" {
		fmt.Println("Commit hash is required.")
		flag.Usage()
		os.Exit(1)
	}

	// Set repository URL and clone directory
	repoURL := "https://github.com/mtumilowicz/helm-workshop"
	cloneDir := "./helm-workshop"

	// Clone the repository if it doesn't exist
	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		fmt.Printf("Cloning repository from %s...\n", repoURL)
		err := cloneRepository(repoURL, cloneDir)
		if err != nil {
			fmt.Printf("Error cloning repository: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Repository already exists at %s, skipping clone.\n", cloneDir)
	}

	// Clean repository and checkout the specific commit
	fmt.Printf("Cleaning repository...\n")
	err := cleanRepository(cloneDir)
	if err != nil {
		fmt.Printf("Error cleaning repository: %v\n", err)
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

// cleanRepository ensures that the git repository is in a clean state.
func cleanRepository(cloneDir string) error {
	cmd := exec.Command("git", "reset", "--hard")
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "clean", "-fd")
	cmd.Dir = cloneDir
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
