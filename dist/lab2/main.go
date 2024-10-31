package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Task represents a unit of work.
type Task struct {
	id         int
	data       string
	retryCount int
}

const maxRetries = 3

// Worker function that processes tasks. If a worker fails, the task will be sent to failChan.
func worker(id int, taskChan <-chan Task, doneChan, failChan chan<- Task) {
	for task := range taskChan {
		logrus.WithFields(logrus.Fields{
			"worker_id": id,
			"task_id":   task.id,
			"task_data": task.data,
			"retries":   task.retryCount,
		}).Info("Worker started processing task")

		// Simulate random failure (30% chance of failure)
		if rand.Float32() < 0.3 {
			logrus.WithFields(logrus.Fields{"worker_id": id, "task_id": task.id}).Error("Worker failed task")
			failChan <- task
			continue
		}

		// Simulate task processing time
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		doneChan <- task
		logrus.WithFields(logrus.Fields{"worker_id": id, "task_id": task.id}).Info("Worker completed task")
	}
	logrus.WithFields(logrus.Fields{"worker_id": id}).Info("Worker channel is closed and worker is shutting down")
}

func loadBalancer(numWorkers int, taskChan <-chan Task, doneChan, failChan chan<- Task) {
	channelMap := make(map[int]chan Task)
	defer func() {
		for _, workerChan := range channelMap {
			close(workerChan)
		}
		logrus.Info("Load balancer closed all worker channels")
	}()

	for id := 0; id < numWorkers; id++ {
		channelMap[id] = make(chan Task)
		go worker(id, channelMap[id], doneChan, failChan)
		logrus.WithFields(logrus.Fields{"worker_id": id}).Info("Worker initialized and ready")
	}

	next_worker_id := 0
	for task := range taskChan {
		workerChan := channelMap[next_worker_id]
		logrus.WithFields(logrus.Fields{"task_id": task.id, "worker_id": next_worker_id}).Info("Task assigned to worker")
		workerChan <- task
		next_worker_id = (next_worker_id + 1) % numWorkers
	}
	logrus.Info("All tasks have been processed. Shutting down load balancer.")
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	// Define a set of tasks to be executed
	tasks := []Task{
		{id: 1, data: "Task 1"},
		{id: 2, data: "Task 2"},
		{id: 3, data: "Task 3"},
		{id: 4, data: "Task 4"},
		{id: 5, data: "Task 5"},
	}

	// Channels for task distribution and failure handling
	taskChan := make(chan Task, len(tasks))
	failChan := make(chan Task, len(tasks))
	doneChan := make(chan Task, len(tasks))

	// Use a WaitGroup to wait for all tasks to complete
	var wg sync.WaitGroup

	// Distribute tasks to the available workers
	for _, task := range tasks {
		logrus.WithFields(logrus.Fields{"task_id": task.id}).Info("Dispatching task to task channel")
		taskChan <- task
		wg.Add(1)
	}

	go loadBalancer(3, taskChan, doneChan, failChan)

	// Handle failed tasks by redistributing them
	go func() {
		for failedTask := range failChan {
			failedTask.retryCount++
			if failedTask.retryCount > maxRetries {
				logrus.WithFields(logrus.Fields{
					"task_id":    failedTask.id,
					"retryCount": failedTask.retryCount,
				}).Warn("Task failed; max retry limit reached")
				wg.Done()
				continue
			}
			logrus.WithFields(logrus.Fields{
				"task_id":    failedTask.id,
				"retryCount": failedTask.retryCount,
			}).Info("Retrying failed task")
			taskChan <- failedTask // Send back to the task channel to be retried
		}
	}()

	// Collect completed tasks and close channels when done
	go func() {
		for range doneChan {
			wg.Done()
			logrus.Info("Task completed and marked done")
		}
	}()

	wg.Wait()

	defer func() {
		close(taskChan)
		close(failChan)
		close(doneChan)
		logrus.Info("All tasks have been processed. Shutting down main.")
	}()
}
