# ssh-agent auto-start for fish (conf.d)
# Ensures a single ssh-agent per login session and preserves ssh-add across terminals
# - Reuses environment from ~/.cache/ssh-agent.env when available
# - Starts a new agent only if necessary (and only when the socket is invalid)
# - Adds common default keys non-interactively if present

# Only for interactive shells
if not status is-interactive
    exit
end

# Dependencies check
if not command -q ssh-agent
    # ssh-agent not available; nothing to do
    exit
end
if not command -q ssh-add
    # ssh-add not available; nothing to do
    exit
end

set -l envfile ~/.cache/ssh-agent.env

function __ssh_agent__save_env --description 'Save SSH agent environment for reuse'
    mkdir -p (dirname $envfile) ^/dev/null
    printf 'set -gx SSH_AUTH_SOCK %s\nset -gx SSH_AGENT_PID %s\n' $SSH_AUTH_SOCK $SSH_AGENT_PID > $envfile
end

function __ssh_agent__start --description 'Start ssh-agent and export environment'
    set -l out (ssh-agent -s)
    # Parse SSH_AUTH_SOCK and SSH_AGENT_PID from sh-style output
    set -l sock (string match -r 'SSH_AUTH_SOCK=([^;]+);' $out | string replace -r '.*SSH_AUTH_SOCK=([^;]+);.*' '$1')
    set -l pid  (string match -r 'SSH_AGENT_PID=([0-9]+);' $out | string replace -r '.*SSH_AGENT_PID=([0-9]+);.*' '$1')

    if test -n "$sock"; and test -n "$pid"
        set -gx SSH_AUTH_SOCK $sock
        set -gx SSH_AGENT_PID $pid
        __ssh_agent__save_env
        return 0
    end
    return 1
end

# Try to ensure agent env is valid
set -l have_valid 0

# Case 1: Current environment is already valid
if set -q SSH_AUTH_SOCK; and test -S "$SSH_AUTH_SOCK"
    set have_valid 1
else
    # Case 2: Try to source a previously saved environment
    if test -f $envfile
        source $envfile ^/dev/null
        if set -q SSH_AUTH_SOCK; and test -S "$SSH_AUTH_SOCK"
            set have_valid 1
        end
    end
end

# Case 3: If still invalid, start a new agent (prefer on login shells)
if test $have_valid -eq 0
    if status is-login
        __ssh_agent__start >/dev/null 2>&1
        if set -q SSH_AUTH_SOCK; and test -S "$SSH_AUTH_SOCK"
            set have_valid 1
        end
    end
end

# If we have a valid agent, ensure default keys are added (quietly)
if test $have_valid -eq 1
    # If no identities are loaded, attempt to add common defaults (if present)
    ssh-add -l >/dev/null 2>&1
    if test $status -ne 0
        set -l keys id_ed25519 id_rsa id_ecdsa id_dsa
        for k in $keys
            if test -f ~/.ssh/$k
                ssh-add ~/.ssh/$k >/dev/null 2>&1
            end
        end
    end
end
