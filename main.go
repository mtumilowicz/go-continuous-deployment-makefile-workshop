package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Define flags
	action := flag.String("action", "", "Action to perform (clone, checkout, clean, test, build, upgrade)")
	commitHash := flag.String("commit-hash", "", "Commit hash to use")
	imageName := flag.String("image-name", "", "Name of the Docker image")
	releaseName := flag.String("release-name", "", "Helm release name")
	repoURL := flag.String("repo-url", "", "URL of the Git repository to clone")
	chartDir := flag.String("chart-dir", "", "Directory of the Helm chart")
	namespace := flag.String("namespace", "default", "Kubernetes namespace for the Helm release")
	imageVersion := flag.String("image-version", "", "Image version for Helm release")

	// Parse the flags
	flag.Parse()

	// Validate action
	if *action == "" {
		fmt.Println("ACTION is required.")
		os.Exit(1)
	}

	// Get and validate clone directory based on image name
	cloneDir := getCloneDir(*imageName)

	// Perform the action based on parameters
	switch *action {
	case "clone":
		cloneRepository(*repoURL, cloneDir)
	case "checkout":
		checkoutCommit(*commitHash, cloneDir)
	case "clean":
		cleanBuild(cloneDir)
	case "test":
		runTests(cloneDir)
	case "build":
		buildDockerImage(*imageVersion, *imageName, cloneDir)
	case "upgrade":
		upgradeHelmChart(*imageVersion, *releaseName, *chartDir, *namespace, cloneDir)
	default:
		fmt.Println("Unknown action:", *action)
		os.Exit(1)
	}
}

// getCloneDir validates the image name and returns the clone directory.
func getCloneDir(imageName string) string {
	if imageName == "" {
		fmt.Println("IMAGE_NAME is required.")
		os.Exit(1)
	}
	return "./" + imageName
}

// cloneRepository clones the Git repository to the specified directory.
func cloneRepository(repoURL, cloneDir string) {
	if repoURL == "" {
		fmt.Println("REPO_URL is required for clone action.")
		os.Exit(1)
	}
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Repository cloned successfully.")
}

// checkoutCommit checks out the specified commit hash in the working directory.
func checkoutCommit(commitHash, cloneDir string) {
	if commitHash == "" {
		fmt.Println("COMMIT_HASH is required for checkout action.")
		os.Exit(1)
	}
	cmd := exec.Command("git", "checkout", commitHash)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error checking out commit: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Checked out commit successfully.")
}

// cleanBuild removes build artifacts from the working directory.
func cleanBuild(cloneDir string) {
	cmd := exec.Command("./gradlew", "clean")
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error cleaning build artifacts: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build artifacts cleaned successfully.")
}

// runTests runs the tests using Gradle in the working directory.
func runTests(cloneDir string) {
	cmd := exec.Command("./gradlew", "test")
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Tests ran successfully.")
}

// buildDockerImage builds the Docker image using Gradle with the given commit hash as the tag.
func buildDockerImage(commitHash, imageName, cloneDir string) {
	if commitHash == "" {
		fmt.Println("COMMIT_HASH is required for build action.")
		os.Exit(1)
	}
	tag := fmt.Sprintf("%s:%s", imageName, commitHash)
	cmd := exec.Command("./gradlew", "bootBuildImage", "--imageName="+tag)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Docker image built successfully.")
}

// upgradeHelmChart upgrades the Helm chart with the given image tag, chart directory, and release name.
func upgradeHelmChart(imageVersion, releaseName, chartDir, namespace, cloneDir string) {
	if imageVersion == "" || releaseName == "" || chartDir == "" {
		fmt.Println("IMAGE_VERSION, RELEASE_NAME, and CHART_DIR are required for upgrade action.")
		os.Exit(1)
	}
	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartDir, "--set", fmt.Sprintf("deployment.image.version=%s", imageVersion), "--namespace", namespace)
	cmd.Dir = cloneDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error upgrading Helm chart: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Helm chart upgraded successfully.")
}
