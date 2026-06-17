#!/bin/bash
# Start infrastructure (detached, exits immediately)
make docker

# Split top pane vertically (left/right)
tmux split-window -h

# Split horizontally (top/bottom)
tmux select-pane -t 0
tmux split-window -v

# Split bottom pane vertically (left/right)
tmux select-pane -t 1
tmux split-window -v

# Run commands in panes
tmux select-pane -t 0
tmux send-keys 'make mobile' C-m

tmux select-pane -t 1
tmux send-keys 'make server' C-m

# Attach to session
tmux attach -t lazyfood-dev
