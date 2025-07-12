#!/usr/bin/env python3
"""
Enhanced Tomato Timer (Pomodoro) for Waybar
- 25 minute work session
- 1 minute automatic transition wait
- 15 minute break countdown (burnt orange)
- Return to 25 minute ready state
"""
import json
import os
import sys
import time
from datetime import datetime, timedelta

TIMER_FILE = "/tmp/waybar-tomato-timer.data"
STATE_FILE = "/tmp/waybar-tomato-timer.state"

class EnhancedTomatoTimer:
    def __init__(self):
        self.work_minutes = 25
        self.break_minutes = 15
        self.transition_seconds = 60  # 1 minute transition wait
        self.load_timer_data()
        self.check_for_commands()

    def load_timer_data(self):
        """Load timer data from file"""
        try:
            if os.path.exists(TIMER_FILE):
                with open(TIMER_FILE, 'r') as f:
                    data = json.loads(f.read())
                    self.mode = data.get('mode', 'idle')  # idle, work, transition, break
                    self.is_running = data.get('is_running', False)
                    self.original_duration = data.get('original_duration', self.work_minutes * 60)

                    # Only set end_time if we have a valid one and timer was running
                    if data.get('end_time'):
                        try:
                            self.end_time = datetime.fromisoformat(data['end_time'])
                        except:
                            self.reset_timer()
                    else:
                        # No end_time - fresh start
                        if hasattr(self, 'end_time'):
                            delattr(self, 'end_time')
            else:
                self.reset_timer()
        except:
            self.reset_timer()

    def save_timer_data(self):
        """Save timer data to file"""
        data = {
            'mode': getattr(self, 'mode', 'idle'),
            'end_time': self.end_time.isoformat() if hasattr(self, 'end_time') else None,
            'is_running': getattr(self, 'is_running', False),
            'original_duration': getattr(self, 'original_duration', self.work_minutes * 60)
        }
        with open(TIMER_FILE, 'w') as f:
            f.write(json.dumps(data))

    def check_for_commands(self):
        """Check for click commands from the click handler"""
        if os.path.exists(STATE_FILE):
            try:
                with open(STATE_FILE, 'r') as f:
                    cmd = json.loads(f.read())

                action = cmd.get('action')
                if action == 'toggle':
                    self.start_pause()
                elif action == 'reset':
                    self.reset_timer()

                # Remove the command file
                os.remove(STATE_FILE)
            except:
                # Invalid command file, remove it
                try:
                    os.remove(STATE_FILE)
                except:
                    pass

    def reset_timer(self):
        """Reset to 25 minute work mode and stop"""
        self.mode = 'idle'
        self.is_running = False
        self.original_duration = self.work_minutes * 60
        if hasattr(self, 'end_time'):
            delattr(self, 'end_time')
        self.save_timer_data()

    def start_pause(self):
        """Start or pause the timer"""
        now = datetime.now()

        # Can only manually control work sessions, not break sessions
        if self.mode in ['transition', 'break']:
            return  # Do nothing during automatic break sequence

        if self.mode == 'idle':
            # Start new 25 minute work session
            self.mode = 'work'
            self.is_running = True
            self.end_time = now + timedelta(minutes=self.work_minutes)
            self.original_duration = self.work_minutes * 60
        elif self.mode == 'work':
            # Toggle pause/resume work session
            if self.is_running:
                # Pause - store remaining time
                remaining = (self.end_time - now).total_seconds()
                self.remaining_seconds = max(0, remaining)
                self.is_running = False
            else:
                # Resume - set new end time based on remaining time
                self.end_time = now + timedelta(seconds=self.remaining_seconds)
                self.is_running = True

        self.save_timer_data()

    def start_transition(self):
        """Start the transition period after work completion"""
        self.mode = 'transition'
        self.is_running = True
        self.end_time = datetime.now() + timedelta(seconds=self.transition_seconds)
        self.save_timer_data()

    def start_break(self):
        """Start the break period"""
        self.mode = 'break'
        self.is_running = True
        self.end_time = datetime.now() + timedelta(minutes=self.break_minutes)
        self.save_timer_data()

    def get_display(self):
        """Get current display state for waybar"""
        now = datetime.now()

        # Handle mode transitions and timing
        if hasattr(self, 'end_time') and self.is_running:
            remaining = (self.end_time - now).total_seconds()

            if remaining <= 0:
                # Timer finished - handle mode transitions
                if self.mode == 'work':
                    # Work finished, start transition period
                    self.start_transition()
                    remaining = self.transition_seconds
                elif self.mode == 'transition':
                    # Transition finished, start break
                    self.start_break()
                    remaining = self.break_minutes * 60
                elif self.mode == 'break':
                    # Break finished, reset to idle
                    self.reset_timer()
                    return self.get_idle_display()

        # Display based on current mode
        if self.mode == 'idle' or not hasattr(self, 'end_time'):
            return self.get_idle_display()
        elif self.mode == 'work':
            return self.get_work_display()
        elif self.mode == 'transition':
            return self.get_transition_display()
        elif self.mode == 'break':
            return self.get_break_display()
        else:
            return self.get_idle_display()

    def get_idle_display(self):
        """Display for idle/ready state"""
        return {
            "text": "󰌾 25:00",
            "tooltip": "Pomodoro Timer (Ready)\nClick: Start 25min work session\nDouble-click: Reset",
            "class": "idle"
        }

    def get_work_display(self):
        """Display for work session"""
        now = datetime.now()
        remaining = (self.end_time - now).total_seconds()

        if not self.is_running:
            # Paused work session
            minutes = int(self.remaining_seconds // 60)
            seconds = int(self.remaining_seconds % 60)
            return {
                "text": f"󰏤 {minutes:02d}:{seconds:02d}",
                "tooltip": "Work Session (Paused)\nClick: Resume\nDouble-click: Reset",
                "class": "paused"
            }

        # Running work session
        if remaining <= 0:
            # Work session completed
            blink = int(time.time()) % 2 == 0
            icon = "󰂚" if blink else "󰌾"
            return {
                "text": f"{icon} Work Done!",
                "tooltip": "Work Session Complete!\nStarting break in a moment...",
                "class": "finished"
            }

        minutes = int(remaining // 60)
        seconds = int(remaining % 60)
        return {
            "text": f"󰔛 {minutes:02d}:{seconds:02d}",
            "tooltip": f"Work Session ({self.work_minutes}min)\nClick: Pause\nDouble-click: Reset",
            "class": "work"
        }

    def get_transition_display(self):
        """Display for transition period between work and break"""
        now = datetime.now()
        remaining = (self.end_time - now).total_seconds()

        if remaining <= 0:
            return {
                "text": "󰌾 Starting Break...",
                "tooltip": "Starting your break now...",
                "class": "transition"
            }

        seconds = int(remaining)
        blink = int(time.time()) % 2 == 0
        icon = "󰂚" if blink else "󰌾"
        return {
            "text": f"{icon} Break in {seconds}s",
            "tooltip": "Work complete! Break starting soon...",
            "class": "transition"
        }

    def get_break_display(self):
        """Display for break period"""
        now = datetime.now()
        remaining = (self.end_time - now).total_seconds()

        if remaining <= 0:
            return {
                "text": "󰌾 Break Over!",
                "tooltip": "Break finished! Ready for next session.",
                "class": "break-finished"
            }

        minutes = int(remaining // 60)
        seconds = int(remaining % 60)
        blink = int(time.time() * 2) % 2 == 0  # Faster blink for break
        icon = "☕" if blink else "󰌾"
        return {
            "text": f"{icon} {minutes:02d}:{seconds:02d}",
            "tooltip": f"Break Time ({self.break_minutes}min)\nEnjoy your break!",
            "class": "break"
        }

# Create timer instance and get display
timer = EnhancedTomatoTimer()
print(json.dumps(timer.get_display()))
