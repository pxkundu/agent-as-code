#!/usr/bin/env python3
"""
Test Enhanced Features
=====================

Simple test script to verify that the enhanced LLM features
are properly accessible in the Python package.
"""

import sys
import os

# Add the current directory to Python path
sys.path.insert(0, os.path.dirname(__file__))

def test_import():
    """Test that the package can be imported."""
    try:
        import agent_as_code
        print(f"✅ Package imported successfully: {agent_as_code.__version__}")
        return True
    except ImportError as e:
        print(f"❌ Failed to import package: {e}")
        return False

def test_agent_cli():
    """Test that AgentCLI can be instantiated."""
    try:
        from agent_as_code import AgentCLI
        cli = AgentCLI()
        print("✅ AgentCLI instantiated successfully")
        return cli
    except Exception as e:
        print(f"❌ Failed to instantiate AgentCLI: {e}")
        return None

def test_enhanced_methods(cli):
    """Test that all enhanced methods are available."""
    if not cli:
        return False
    
    enhanced_methods = [
        'create_agent',
        'optimize_model', 
        'benchmark_models',
        'deploy_agent',
        'analyze_model',
        'list_models',
        'pull_model',
        'test_model',
        'remove_model'
    ]
    
    print("\n🔍 Testing enhanced methods availability:")
    all_available = True
    
    for method in enhanced_methods:
        if hasattr(cli, method):
            print(f"  ✅ {method}")
        else:
            print(f"  ❌ {method}")
            all_available = False
    
    return all_available

def test_method_signatures(cli):
    """Test that methods have proper signatures."""
    if not cli:
        return False
    
    print("\n🔍 Testing method signatures:")
    
    # Test create_agent method
    try:
        import inspect
        sig = inspect.signature(cli.create_agent)
        params = list(sig.parameters.keys())
        # Bound methods don't include 'self' in the signature
        expected_params = ['use_case', 'model', 'optimize', 'test']
        
        if params == expected_params:
            print("  ✅ create_agent signature correct")
        else:
            print(f"  ❌ create_agent signature mismatch: {params} vs {expected_params}")
            return False
            
    except Exception as e:
        print(f"  ❌ Failed to inspect create_agent signature: {e}")
        return False
    
    return True

def main():
    """Run all tests."""
    print("🧠 Testing Enhanced LLM Features")
    print("=" * 40)
    
    # Test 1: Import
    if not test_import():
        print("\n❌ Import test failed")
        return 1
    
    # Test 2: AgentCLI instantiation
    cli = test_agent_cli()
    if not cli:
        print("\n❌ AgentCLI test failed")
        return 1
    
    # Test 3: Enhanced methods availability
    if not test_enhanced_methods(cli):
        print("\n❌ Enhanced methods test failed")
        return 1
    
    # Test 4: Method signatures
    if not test_method_signatures(cli):
        print("\n❌ Method signatures test failed")
        return 1
    
    print("\n🎉 All tests passed!")
    print("\n💡 The Python package is ready with enhanced LLM features.")
    print("   You can now use:")
    print("   - cli.create_agent('chatbot')")
    print("   - cli.optimize_model('llama2', 'chatbot')")
    print("   - cli.benchmark_models(['chatbot', 'analysis'])")
    print("   - cli.deploy_agent('my-agent')")
    print("   - cli.analyze_model('llama2', detailed=True)")
    
    return 0

if __name__ == "__main__":
    exit(main())
