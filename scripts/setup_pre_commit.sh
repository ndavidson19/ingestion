#!/bin/bash

set -e

# Install pre-commit
pip install pre-commit

# Install the pre-commit hooks
pre-commit install

echo "Pre-commit hooks have been installed successfully."
