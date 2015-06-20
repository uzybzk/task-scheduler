package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Task struct {
    ID       int
    Name     string
    Schedule string
    Command  string
    NextRun  time.Time
    Enabled  bool
}

type Scheduler struct {
    tasks   []Task
    running bool
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        tasks:   make([]Task, 0),
        running: false,
    }
}

func (s *Scheduler) AddTask(name, schedule, command string) {
    nextRun := calculateNextRun(schedule)
    task := Task{
        ID:       len(s.tasks) + 1,
        Name:     name,
        Schedule: schedule,
        Command:  command,
        NextRun:  nextRun,
        Enabled:  true,
    }
    s.tasks = append(s.tasks, task)
    fmt.Printf("Added task: %s (ID: %d)\n", name, task.ID)
}

func (s *Scheduler) Start() {
    s.running = true
    fmt.Println("Task scheduler started")
    
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for s.running {
        select {
        case <-ticker.C:
            s.checkTasks()
        }
    }
}

func (s *Scheduler) Stop() {
    s.running = false
    fmt.Println("Task scheduler stopped")
}

func (s *Scheduler) checkTasks() {
    now := time.Now()
    
    for i := range s.tasks {
        task := &s.tasks[i]
        
        if task.Enabled && now.After(task.NextRun) {
            fmt.Printf("Executing task: %s\n", task.Name)
            executeTask(task)
            
            // Calculate next run time
            task.NextRun = calculateNextRun(task.Schedule)
            fmt.Printf("Next run for %s: %s\n", task.Name, task.NextRun.Format("2006-01-02 15:04:05"))
        }
    }
}

func (s *Scheduler) ListTasks() {
    fmt.Println("Scheduled Tasks:")
    fmt.Println("================")
    for _, task := range s.tasks {
        status := "Enabled"
        if !task.Enabled {
            status = "Disabled"
        }
        fmt.Printf("ID: %d | Name: %s | Schedule: %s | Status: %s | Next: %s\n",
            task.ID, task.Name, task.Schedule, status, task.NextRun.Format("2006-01-02 15:04:05"))
    }
}

func executeTask(task *Task) {
    // Simple task execution simulation
    fmt.Printf("Running command: %s\n", task.Command)
    
    // In a real implementation, this would execute the actual command
    // For demo purposes, we just log it
    time.Sleep(100 * time.Millisecond) // Simulate work
    fmt.Printf("Task %s completed\n", task.Name)
}

func calculateNextRun(schedule string) time.Time {
    now := time.Now()
    
    switch schedule {
    case "@hourly":
        return now.Add(1 * time.Hour)
    case "@daily":
        return now.AddDate(0, 0, 1)
    case "@weekly":
        return now.AddDate(0, 0, 7)
    case "@monthly":
        return now.AddDate(0, 1, 0)
    default:
        // Default to hourly if unknown
        return now.Add(1 * time.Hour)
    }
}

func main() {
    scheduler := NewScheduler()
    
    // Add some sample tasks
    scheduler.AddTask("Backup Database", "@daily", "pg_dump mydb > backup.sql")
    scheduler.AddTask("Clean Logs", "@weekly", "find /var/log -name '*.log' -mtime +7 -delete")
    scheduler.AddTask("System Health Check", "@hourly", "systemctl status important-service")
    
    // List current tasks
    scheduler.ListTasks()
    
    // Set up graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        fmt.Println("\nReceived interrupt signal, shutting down...")
        scheduler.Stop()
        os.Exit(0)
    }()
    
    fmt.Println("Starting task scheduler...")
    fmt.Println("Press Ctrl+C to stop")
    
    scheduler.Start()
}