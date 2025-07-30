#!/usr/bin/env python3
"""Simple Pomodoro Timer for Waybar"""
import json
import os
from datetime import datetime, timedelta
from pathlib import Path

TIMER_FILE = "/tmp/waybar-tomato.json"
STATE_FILE = "/tmp/waybar-tomato-timer.state"
CONFIG_FILE = Path.home() / ".config/archriot/archriot.conf"

class SimplePomodoro:
    def __init__(self):
        self.enabled, self.work_minutes = self.load_config()
        self.break_minutes = 5
        self.load_and_process()

    def load_config(self):
        """Load Pomodoro config from ArchRiot config file"""
        try:
            import configparser
            config = configparser.ConfigParser()
            config.read(CONFIG_FILE)
            enabled = config.getboolean('pomodoro', 'enabled', fallback=True)
            duration = config.getint('pomodoro', 'duration', fallback=25)
            return enabled, duration
        except:
            return True, 25  # Default fallback

    def load_and_process(self):
        """Load state and process commands"""
        # Load timer state
        try:
            if os.path.exists(TIMER_FILE):
                with open(TIMER_FILE, 'r') as f:
                    data = json.load(f)
                self.mode = data.get('mode', 'idle')
                self.running = data.get('running', False)
                self.end_time = datetime.fromisoformat(data['end_time']) if data.get('end_time') else None
                self.paused_remaining = data.get('paused_remaining', None)
            else:
                self.reset_state()
        except:
            self.reset_state()

        # Process click commands - but only if enabled
        if self.enabled and os.path.exists(STATE_FILE):
            try:
                with open(STATE_FILE, 'r') as f:
                    action = json.load(f).get('action')
                if action == 'toggle':
                    self.toggle()
                elif action == 'reset':
                    self.reset_state()
                    self.save_state()
                os.remove(STATE_FILE)
            except:
                pass
        elif not self.enabled and os.path.exists(STATE_FILE):
            # Remove state file but don't process when disabled
            try:
                os.remove(STATE_FILE)
            except:
                pass

    def reset_state(self):
        """Reset to idle state"""
        self.mode = 'idle'
        self.running = False
        self.end_time = None
        self.paused_remaining = None

    def save_state(self):
        """Save timer state"""
        data = {
            'mode': self.mode,
            'running': self.running,
            'end_time': self.end_time.isoformat() if self.end_time else None,
            'paused_remaining': self.paused_remaining
        }
        with open(TIMER_FILE, 'w') as f:
            json.dump(data, f)

    def toggle(self):
        """Start/pause timer"""
        if self.mode == 'idle':
            self.mode = 'work'
            self.running = True
            self.end_time = datetime.now() + timedelta(minutes=self.work_minutes)
        elif self.running:
            # Pause - save remaining time
            if self.end_time:
                remaining = (self.end_time - datetime.now()).total_seconds()
                self.paused_remaining = max(0, remaining)
            self.running = False
        else:
            # Resume - set new end time based on saved remaining time
            if self.paused_remaining:
                self.end_time = datetime.now() + timedelta(seconds=self.paused_remaining)
                self.paused_remaining = None
            self.running = True
        self.save_state()

    def get_display(self):
        """Get waybar display"""
        # If disabled, show grayed out version
        if not self.enabled:
            return {"text": "󰌾 --:--", "tooltip": "Pomodoro Timer - Disabled", "class": "disabled"}

        if self.mode == 'idle':
            return {"text": f"󰌾 {self.work_minutes:02d}:00", "tooltip": "Pomodoro Timer - Click to start", "class": "idle"}

        if not self.running and self.paused_remaining is not None:
            remaining = self.paused_remaining
        else:
            remaining = (self.end_time - datetime.now()).total_seconds() if self.end_time else 0

        if remaining <= 0:
            if self.mode == 'work':
                self.mode = 'break'
                self.end_time = datetime.now() + timedelta(minutes=self.break_minutes)
                self.save_state()
                return {"text": "󰌾 Break!", "tooltip": "Work done! Take a 5 minute break", "class": "break"}
            else:
                self.reset_state()
                self.save_state()
                return {"text": "󰌾 Ready", "tooltip": "Break over! Ready for next session", "class": "ready"}

        minutes, seconds = int(remaining // 60), int(remaining % 60)
        icon = "󰔛" if self.mode == 'work' and self.running else "󰏤" if not self.running else "☕"
        status = "Work" if self.mode == 'work' and self.running else "Break" if self.mode == 'break' else "Paused"

        return {
            "text": f"{icon} {minutes:02d}:{seconds:02d}",
            "tooltip": f"{status} - {minutes}:{seconds:02d} remaining",
            "class": self.mode if self.running else "paused"
        }

if __name__ == "__main__":
    timer = SimplePomodoro()
    print(json.dumps(timer.get_display()))
