"""
Agent as Code CLI Wrapper
=========================

Python wrapper around the Go-based Agent as Code CLI binary.
This module provides seamless integration between Python and the high-performance
Go implementation while maintaining familiar Python development patterns.
"""

import os
import sys
import platform
import subprocess
import shutil
from pathlib import Path
from typing import List, Optional, Union


class AgentCLI:
    """
    Python wrapper for the Agent as Code CLI.
    
    This class provides a Python interface to the Go-based CLI binary,
    automatically detecting the correct binary for the current platform
    and providing both direct execution and Python API methods.
    """
    
    def __init__(self, binary_path: Optional[str] = None):
        """
        Initialize the Agent CLI wrapper.
        
        Args:
            binary_path: Optional path to the agent binary. If not provided,
                        will auto-detect based on the current platform.
        """
        self.binary_path = binary_path or self._detect_binary()
        self._validate_binary()
    
    def _detect_binary(self) -> Path:
        """
        Auto-detect the correct binary for the current platform.
        
        Returns:
            Path to the appropriate binary for the current platform.
            
        Raises:
            RuntimeError: If the platform is unsupported or binary not found.
        """
        system = platform.system().lower()
        machine = platform.machine().lower()
        
        # Normalize architecture names
        arch_map = {
            'x86_64': 'amd64',
            'amd64': 'amd64',
            'arm64': 'arm64',
            'aarch64': 'arm64',
        }
        
        arch = arch_map.get(machine)
        if not arch:
            raise RuntimeError(f"Unsupported architecture: {machine}")
        
        # Map platform to binary name
        platform_map = {
            'linux': f'agent-linux-{arch}',
            'darwin': f'agent-darwin-{arch}',
            'windows': f'agent-windows-{arch}.exe',
        }
        
        binary_name = platform_map.get(system)
        if not binary_name:
            raise RuntimeError(f"Unsupported platform: {system}")
        
        # Find binary in package
        package_dir = Path(__file__).parent
        binary_path = package_dir / "bin" / binary_name
        
        if not binary_path.exists():
            raise RuntimeError(
                f"Binary not found: {binary_path}\n"
                f"Platform: {system}-{arch}\n"
                f"Expected binary: {binary_name}\n"
                f"This may indicate a corrupted installation. "
                f"Try reinstalling: pip install --force-reinstall agent-as-code"
            )
        
        return binary_path
    
    def _validate_binary(self):
        """
        Validate that the binary exists and is executable.
        
        Raises:
            RuntimeError: If the binary is not found or not executable.
        """
        if not self.binary_path.exists():
            raise RuntimeError(f"Agent binary not found: {self.binary_path}")
        
        # Make executable on Unix systems
        if platform.system() != 'Windows':
            self.binary_path.chmod(0o755)
        
        # Test binary execution
        try:
            result = subprocess.run(
                [str(self.binary_path), "--version"],
                capture_output=True,
                text=True,
                timeout=10
            )
            if result.returncode != 0:
                raise RuntimeError(f"Binary test failed: {result.stderr}")
        except subprocess.TimeoutExpired:
            raise RuntimeError("Binary test timed out")
        except Exception as e:
            raise RuntimeError(f"Binary validation failed: {e}")
    
    def execute(self, args: Optional[List[str]] = None) -> int:
        """
        Execute the agent CLI with the given arguments.
        
        Args:
            args: List of command line arguments. If None, uses sys.argv[1:].
            
        Returns:
            Exit code from the CLI execution.
        """
        if args is None:
            args = sys.argv[1:]
        
        try:
            # Execute the binary
            result = subprocess.run(
                [str(self.binary_path)] + args,
                check=False
            )
            return result.returncode
            
        except KeyboardInterrupt:
            return 130  # Standard exit code for Ctrl+C
        except Exception as e:
            print(f"Error executing agent: {e}", file=sys.stderr)
            return 1
    
    def run_command(self, args: List[str], capture_output: bool = False, 
                   timeout: Optional[float] = None) -> subprocess.CompletedProcess:
        """
        Run a command and return the result.
        
        Args:
            args: Command arguments to pass to the CLI.
            capture_output: Whether to capture stdout/stderr.
            timeout: Timeout in seconds for the command.
            
        Returns:
            CompletedProcess instance with the result.
        """
        try:
            return subprocess.run(
                [str(self.binary_path)] + args,
                capture_output=capture_output,
                text=True,
                timeout=timeout,
                check=False
            )
        except subprocess.TimeoutExpired as e:
            raise RuntimeError(f"Command timed out after {timeout}s: {' '.join(args)}") from e
    
    def get_version(self) -> Optional[str]:
        """
        Get the version of the agent CLI binary.
        
        Returns:
            Version string or None if unable to determine.
        """
        try:
            result = self.run_command(["--version"], capture_output=True, timeout=5)
            if result.returncode == 0:
                return result.stdout.strip()
        except Exception:
            pass
        return None
    
    def init(self, name: str, template: Optional[str] = None, 
             runtime: Optional[str] = None, model: Optional[str] = None) -> bool:
        """
        Initialize a new agent project.
        
        Args:
            name: Name of the agent project.
            template: Template to use (chatbot, sentiment, etc.).
            runtime: Runtime environment (python, nodejs, go).
            model: Model to use (openai/gpt-4, local/llama2, etc.).
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["init", name]
        
        if template:
            args.extend(["--template", template])
        if runtime:
            args.extend(["--runtime", runtime])
        if model:
            args.extend(["--model", model])
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def build(self, path: str = ".", tag: Optional[str] = None, 
              no_cache: bool = False, push: bool = False) -> bool:
        """
        Build an agent from agent.yaml.
        
        Args:
            path: Path to the agent project directory.
            tag: Tag for the built image.
            no_cache: Whether to disable build cache.
            push: Whether to push after building.
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["build"]
        
        if tag:
            args.extend(["-t", tag])
        if no_cache:
            args.append("--no-cache")
        if push:
            args.append("--push")
        
        args.append(path)
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def run(self, image: str, port: Optional[str] = None, 
            env: Optional[List[str]] = None, detach: bool = False,
            name: Optional[str] = None) -> bool:
        """
        Run an agent container.
        
        Args:
            image: Image name to run.
            port: Port mapping (e.g., "8080:8080").
            env: List of environment variables.
            detach: Whether to run in background.
            name: Container name.
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["run"]
        
        if port:
            args.extend(["-p", port])
        if env:
            for e in env:
                args.extend(["-e", e])
        if detach:
            args.append("-d")
        if name:
            args.extend(["--name", name])
        
        args.append(image)
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def push(self, image: str, registry: Optional[str] = None) -> bool:
        """
        Push an agent to a registry.
        
        Args:
            image: Image name to push.
            registry: Registry to push to.
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["push"]
        
        if registry:
            args.extend(["--registry", registry])
        
        args.append(image)
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def pull(self, image: str, registry: Optional[str] = None, 
             quiet: bool = False) -> bool:
        """
        Pull an agent from a registry.
        
        Args:
            image: Image name to pull.
            registry: Registry to pull from.
            quiet: Whether to suppress verbose output.
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["pull"]
        
        if registry:
            args.extend(["--registry", registry])
        if quiet:
            args.append("-q")
        
        args.append(image)
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def images(self, quiet: bool = False, all_images: bool = False) -> Union[bool, List[str]]:
        """
        List agent images.
        
        Args:
            quiet: Whether to only show image IDs.
            all_images: Whether to show all images.
            
        Returns:
            True if successful (when not quiet), or list of image IDs (when quiet).
        """
        args = ["images"]
        
        if quiet:
            args.append("-q")
        if all_images:
            args.append("-a")
        
        result = self.run_command(args, capture_output=quiet)
        
        if quiet and result.returncode == 0:
            return result.stdout.strip().split('\n') if result.stdout.strip() else []
        
        return result.returncode == 0

    # Enhanced LLM Commands
    
    def create_agent(self, use_case: str, model: Optional[str] = None, 
                     optimize: bool = False, test: bool = False) -> bool:
        """
        Create an intelligent, fully functional AI agent optimized for a specific use case.
        
        This command uses LLM intelligence to:
        - Generate optimized code based on the use case
        - Create comprehensive test suites
        - Set up proper error handling and logging
        - Configure optimal model parameters
        - Generate deployment configurations
        - Create detailed documentation
        
        Args:
            use_case: The use case for the agent (chatbot, sentiment-analyzer, etc.)
            model: Optional specific model to use
            optimize: Whether to optimize the model for the use case
            test: Whether to run tests after creation
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "create-agent", use_case]
        
        if model:
            args.extend(["--model", model])
        if optimize:
            args.append("--optimize")
        if test:
            args.append("--test")
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def optimize_model(self, model: str, use_case: str) -> bool:
        """
        Optimize a local LLM model for a specific use case.
        
        This command analyzes the model and use case to:
        - Adjust model parameters (temperature, top_p, etc.)
        - Create custom prompts and system messages
        - Optimize context window usage
        - Generate performance benchmarks
        - Create use case specific configurations
        
        Args:
            model: The model to optimize
            use_case: The use case to optimize for
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "optimize", model, use_case]
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def benchmark_models(self, tasks: Optional[List[str]] = None, 
                        output_format: Optional[str] = None) -> bool:
        """
        Run comprehensive benchmarks on all local LLM models.
        
        This command tests models across multiple dimensions:
        - Response time and throughput
        - Memory usage and efficiency
        - Quality assessment for different tasks
        - Cost-benefit analysis
        - Performance recommendations
        
        Args:
            tasks: Optional list of specific tasks to benchmark
            output_format: Optional output format (json, table, etc.)
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "benchmark"]
        
        if tasks:
            args.extend(["--tasks", ",".join(tasks)])
        if output_format:
            args.extend(["--output", output_format])
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def deploy_agent(self, agent_name: str, test_suite: Optional[str] = None,
                     monitor: bool = False) -> bool:
        """
        Deploy and comprehensively test an agent on your local machine.
        
        This command:
        - Builds the agent container
        - Deploys it locally
        - Runs automated tests
        - Validates functionality
        - Provides performance metrics
        - Generates deployment report
        
        Args:
            agent_name: Name of the agent to deploy
            test_suite: Optional test suite to run (basic, comprehensive)
            monitor: Whether to enable monitoring
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "deploy-agent", agent_name]
        
        if test_suite:
            args.extend(["--test-suite", test_suite])
        if monitor:
            args.append("--monitor")
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def analyze_model(self, model: str, detailed: bool = False,
                     capabilities: bool = False) -> bool:
        """
        Analyze a local LLM model's capabilities and limitations.
        
        This command provides deep insights into:
        - Model architecture and parameters
        - Performance characteristics
        - Best use cases and limitations
        - Optimization opportunities
        - Integration recommendations
        
        Args:
            model: The model to analyze
            detailed: Whether to provide detailed analysis
            capabilities: Whether to focus on capabilities
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "analyze", model]
        
        if detailed:
            args.append("--detailed")
        if capabilities:
            args.append("--capabilities")
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def list_models(self, quiet: bool = False) -> Union[bool, List[str]]:
        """
        List available local LLM models.
        
        Args:
            quiet: Whether to only show model names
            
        Returns:
            True if successful (when not quiet), or list of model names (when quiet).
        """
        args = ["llm", "list"]
        
        if quiet:
            args.append("-q")
        
        result = self.run_command(args, capture_output=quiet)
        
        if quiet and result.returncode == 0:
            return result.stdout.strip().split('\n') if result.stdout.strip() else []
        
        return result.returncode == 0
    
    def pull_model(self, model: str, quiet: bool = False) -> bool:
        """
        Pull a model from Ollama.
        
        Args:
            model: The model to pull
            quiet: Whether to suppress verbose output
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "pull", model]
        
        if quiet:
            args.append("-q")
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def test_model(self, model: str, input_text: Optional[str] = None) -> bool:
        """
        Test a local LLM model.
        
        Args:
            model: The model to test
            input_text: Optional input text for testing
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "test", model]
        
        if input_text:
            args.extend(["--input", input_text])
        
        result = self.run_command(args)
        return result.returncode == 0
    
    def remove_model(self, model: str, force: bool = False) -> bool:
        """
        Remove a local LLM model.
        
        Args:
            model: The model to remove
            force: Whether to force removal without confirmation
            
        Returns:
            True if successful, False otherwise.
        """
        args = ["llm", "remove", model]
        
        if force:
            args.append("-f")
        
        result = self.run_command(args)
        return result.returncode == 0


def main():
    """
    Main entry point for the agent CLI command.
    
    This function is called when the 'agent' command is executed from the
    command line after installing the package via pip.
    """
    try:
        cli = AgentCLI()
        exit_code = cli.execute()
        sys.exit(exit_code)
    except RuntimeError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
    except KeyboardInterrupt:
        sys.exit(130)
    except Exception as e:
        print(f"Unexpected error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()