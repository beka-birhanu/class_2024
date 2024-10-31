import numpy as np


class LinearRegressionModel:
    def fit(self, X: np.ndarray, Y: np.ndarray, learning_rate: float, iterations: int) -> tuple[np.ndarray, float]:
        m = Y.size
        self.theta = np.zeros((X.shape[1], 1))  # Initialize weights
        cost_history = []

        for i in range(iterations):
            # Prediction
            y_pred = np.dot(X, self.theta)

            # Reshape Y to be a column vector
            Y = Y.reshape(-1, 1)

            # Cost calculation (Mean Squared Error)
            cost = (1 / (2 * m)) * np.sum(np.square(y_pred - Y))

            # Gradient calculation
            d_theta = (1 / m) * np.dot(X.T, (y_pred - Y))

            # Update weights
            self.theta = self.theta - learning_rate * d_theta

            # Store the cost in history
            cost_history.append(cost)

        return self.theta, cost_history

    def tune(self, X_train: np.ndarray, Y_train: np.ndarray, X_val: np.ndarray, Y_val: np.ndarray, learning_rates: list[float], iterations: int) -> tuple[float, dict]:
        best_learning_rate = None
        best_theta = None
        best_cost = float('inf')
        cost_history = {}

        for lr in learning_rates:
            print(f"Tuning learning rate: {lr}")
            # Train the model with the current learning rate
            theta, _ = self.fit(
                X_train, Y_train, learning_rate=lr, iterations=iterations)

            # Validate the model using the validation data
            val_pred = np.dot(X_val, theta)
            val_cost = (1 / (2 * Y_val.size)) * \
                np.sum(np.square(val_pred - Y_val.reshape(-1, 1)))

            print(f"Validation Cost with learning rate {lr}: {val_cost:.4f}")

            # Check if this is the best learning rate found so far
            if val_cost < best_cost:
                best_cost = val_cost
                best_learning_rate = lr
                best_theta = theta

        print(
            f"Best learning rate: {best_learning_rate} with validation cost: {best_cost:.4f}")

        self.theta = best_theta
        return best_learning_rate, cost_history

    def predict(self, X: np.ndarray) -> np.ndarray:
        """Make predictions using the trained model."""
        return np.dot(X, self.theta)
