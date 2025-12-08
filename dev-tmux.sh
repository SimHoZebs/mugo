#!/bin/bash
# Split top pane vertically (left/right)
tmux split-window -h

# Split horizontally (top/bottom)
tmux select-pane -t 0
tmux split-window -v

# Split bottom pane vertically (left/right)  
tmux select-pane -t 2
tmux split-window -v

# Run commands in panes
tmux select-pane -t 0
tmux send-keys 'make adk-api' C-m

tmux select-pane -t 2
tmux send-keys 'make server' C-m

tmux select-pane -t 1
tmux send-keys 'make emulator && sleep 5 && make mobile' C-m

# Focus on empty pane (bottom-right)
tmux select-pane -t 3

# Attach to session
tmux attach -t lazyfood-dev
