#!/bin/bash
# Script to set up local plugin for Traefik development

PLUGIN_DIR="plugins-local/src/github.com/olepetrik/traefik-geoipstate"

echo "Setting up local plugin directory..."

# Create directory structure
mkdir -p "$PLUGIN_DIR"

# Copy plugin files
cp .traefik.yml go.mod go.sum middleware.go "$PLUGIN_DIR/"

echo "✓ Plugin files copied to $PLUGIN_DIR"
echo ""
echo "Next steps:"
echo "1. Add this to your traefik.yml:"
echo ""
cat << 'YAML'
experimental:
  localPlugins:
    geoipstate:
      moduleName: github.com/olepetrik/traefik-geoipstate
YAML
echo ""
echo "2. Mount the plugins-local directory in your Traefik container"
echo "   volumes:"
echo "     - ./plugins-local:/plugins-local"