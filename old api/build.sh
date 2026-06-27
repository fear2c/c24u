#!/bin/bash

echo "Building CNC for Linux..."
cd cnc
go build -o spectate-cnc
echo "Done! Binary: cnc/spectate-cnc"
echo ""
echo "To run: ./spectate-cnc"
