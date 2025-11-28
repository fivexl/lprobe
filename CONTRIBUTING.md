# Contributing to LProbe

Thank you for your interest in contributing to LProbe! This document provides guidelines for contributing to the project, including how to publish Docker images to your own GitHub Container Registry.

## Publishing Custom Docker Images

This project has been configured to allow contributors to easily publish Docker images to their own GitHub Container Registry (GHCR) space. This is useful when you want to:

- Test your changes in your own environment
- Maintain a fork with custom modifications
- Publish images for your own use or organization

### Quick Setup

To publish Docker images to your own GHCR space, follow these steps:

#### 1. Fork the Repository

First, fork the repository to your GitHub account if you haven't already.

#### 2. Configure Repository Variables

In your forked repository, set up the following repository variables:

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Secrets and variables** → **Actions** → **Variables**
3. Click **New repository variable**
4. Add the following variables:

| Variable | Value | Description |
|----------|-------|-------------|
| `REGISTRY` | `ghcr.io` | The container registry to use (GHCR is recommended) |
| `IMAGE_NAME` | `yourusername/lprobe` | Your custom image name (replace `yourusername` with your GitHub username) |

**Example:**
- `REGISTRY`: `ghcr.io`
- `IMAGE_NAME`: `JohnPreston/lprobe`

#### 3. Trigger a Build

Once the variables are set, you can trigger a Docker image build by:

- **Creating a release tag**: Push a tag matching `v*` (e.g., `v1.0.0`, `v2.1.3`) to trigger a full release
- **Pushing to main branch**: This will create snapshot builds for testing

### Example Workflow

After setting up the variables, when you push a tag like `v1.0.0` to your fork:

1. The GitHub Actions workflow will automatically start
2. Docker images will be built for both `amd64` and `arm64` architectures
3. Images will be pushed to: `ghcr.io/yourusername/lprobe:v1.0.0` and `ghcr.io/yourusername/lprobe:latest`
4. Multi-architecture manifests will be created to support both platforms

### Using Your Custom Images

Once your images are built, you can pull and use them:

```bash
# Pull your custom image
docker pull ghcr.io/yourusername/lprobe:latest

# Use it in your deployments
docker run --rm ghcr.io/yourusername/lprobe:latest -url http://localhost:8080/health
```

### Alternative Registries

While GHCR is the default, you can configure the project to publish to any Docker-compatible registry:

#### Docker Hub

Set your repository variables to:
- `REGISTRY`: `docker.io`
- `IMAGE_NAME`: `yourusername/lprobe`

#### Google Container Registry (GCR)

Set your repository variables to:
- `REGISTRY`: `gcr.io`
- `IMAGE_NAME`: `your-project/lprobe`

#### Private Registries

For private registries, you'll also need to configure authentication:

1. Add your registry credentials as repository **secrets** (not variables):
   - `DOCKER_USERNAME`: Your registry username
   - `DOCKER_PASSWORD`: Your registry password/token

2. Update the workflow to use these secrets for authentication.

### Testing Changes Locally

Before publishing images, you can test your changes locally:

```bash
# Build the binary
go build -o lprobe *.go

# Test the binary
./lprobe -url http://localhost:8080/health

# Build Docker image locally
docker build -t lprobe:test .

# Test the Docker image
docker run --rm lprobe:test -url http://host.docker.internal:8080/health
```

### Submitting Changes Upstream

When you're ready to contribute your changes back to the upstream repository:

1. Ensure your changes follow the project's coding standards
2. Add tests if applicable
3. Update documentation as needed
4. Submit a pull request with a clear description of your changes

The configurable registry setup ensures that:
- The upstream repository continues to work with its original settings
- Your changes don't break existing functionality
- Contributors can easily test and publish their own versions

### Troubleshooting

#### Build Failures

If your builds fail, check:
1. Repository variables are correctly set
2. Your fork has the necessary permissions
3. The GoReleaser configuration is valid

#### Permission Issues

Ensure your GitHub Actions have the necessary permissions:
1. Go to **Settings** → **Actions** → **General**
2. Under **Workflow permissions**, select **Read and write permissions**

#### Image Pull Issues

If you can't pull your images:
1. Check that the images were successfully built and pushed
2. Verify your repository is public (or you have access to private repositories)
3. Ensure you're using the correct image name and tag

## General Contribution Guidelines

### Code Style

- Follow the existing code style and conventions
- Add comments for complex logic
- Ensure your code is properly tested

### Testing

- Run the test suite before submitting changes: `./scripts/test.sh`
- Add new tests for new functionality
- Ensure all tests pass

### Documentation

- Update relevant documentation when making changes
- Add clear comments for complex features
- Keep README and other docs up to date

### Pull Requests

- Create descriptive pull requests
- Link to relevant issues
- Provide clear descriptions of changes
- Be responsive to review comments

Thank you for contributing to LProbe! 🚀