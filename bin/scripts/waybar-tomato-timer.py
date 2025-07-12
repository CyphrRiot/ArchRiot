#!/usr/bin/env python3
"""Simple Pomodoro Timer for Waybar"""
import json
import os
from datetime import datetime, timedelta

TIMER_FILE = "/tmp/waybar-tomato.json"
STATE_FILE = "/tmp/waybar-tomato-timer.state"

class SimplePomodoro:
    def __init__(self):
        self.work_minutes = 25
        self.break_minutes = 5
        self.load_and_process()

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
            else:
                self.reset_state()
        except:
            self.reset_state()

        # Process click commands
        if os.path.exists(STATE_FILE):
            try:
                with open(STATE_FILE, 'r') as f:
                    action = json.load(f).get('action')
                if action == 'toggle':
                    self.toggle()
                elif action == 'reset':
                    self.reset_state()
                os.remove(STATE_FILE)
            except:
                pass

    def reset_state(self):
        """Reset to idle state"""
        self.mode = 'idle'
        self.running = False
        self.end_time = None

    def save_state(self):
        """Save timer state"""
        data = {
            'mode': self.mode,
            'running': self.running,
            'end_time': self.end_time.isoformat() if self.end_time else None
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
            self.running = False
        else:
            self.running = True
        self.save_state()

    def get_display(self):
        """Get waybar display"""
        if self.mode == 'idle':
            return {"text": "🍅 25:00", "tooltip": "Pomodoro Timer - Click to start", "class": "idle"}

        remaining = (self.end_time - datetime.now()).total_seconds() if self.end_time else 0

        if remaining <= 0:
            if self.mode == 'work':
                self.mode = 'break'
                self.end_time = datetime.now() + timedelta(minutes=self.break_minutes)
                self.save_state()
                return {"text": "🍅 Break!", "tooltip": "Work done! Take a 5 minute break", "class": "break"}
            else:
                self.reset_state()
                self.save_state()
                return {"text": "🍅 Ready", "tooltip": "Break over! Ready for next session", "class": "ready"}

        minutes, seconds = int(remaining // 60), int(remaining % 60)
        icon = "🍅" if self.mode == 'work' and self.running else "☕" if self.mode == 'break' else "⏸️"
        status = "Work" if self.mode == 'work' and self.running else "Break" if self.mode == 'break' else "Paused"

        return {
            "text": f"{icon} {minutes:02d}:{seconds:02d}",
            "tooltip": f"{status} - {minutes}:{seconds:02d} remaining",
            "class": self.mode if self.running else "paused"
        }

if __name__ == "__main__":
    timer = SimplePomodoro()
    print(json.dumps(timer.get_display()))
