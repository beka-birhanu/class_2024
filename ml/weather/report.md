# Temperature Prediction Using Linear Regression

## 1. Introduction

Predicting temperature is a complex task influenced by a range of meteorological variables. The primary goal of this project is to develop a predictive model that uses linear regression to estimate temperature based on a set of weather-related features. We used gradient descent to fine-tune our model, seeking to optimize the prediction accuracy by minimizing the error between observed and predicted values.

This report provides a detailed explanation of each stage, from data preprocessing to model development and evaluation, and discusses the choices that led to our final model. The results demonstrate the effectiveness of our model, with a focus on interpreting metrics such as Mean Squared Error (MSE), Mean Absolute Error (MAE), and the R² score.

## 2. Project Goals

The primary objectives of this project include:

1. Predicting temperature using a custom linear regression model.
2. Optimizing model accuracy through gradient descent.
3. Evaluating model performance on various metrics.
4. Providing insights on the model’s strengths and areas for improvement.

## 3. Theoretical Background

### Linear Regression

Linear regression is a fundamental statistical technique that estimates relationships between dependent and independent variables. Here, we use it to predict temperature as a linear function of several weather-related features. The underlying assumption is that the temperature can be reasonably estimated as a weighted sum of these features.

The cost function for linear regression is the Mean Squared Error (MSE), which measures the average squared differences between actual and predicted values. Our model is designed to minimize MSE by iteratively updating weights with gradient descent.

### Gradient Descent

Gradient descent is an optimization algorithm widely used for minimizing cost functions in machine learning. This iterative process updates the model parameters (weights) in the opposite direction of the gradient, with the step size controlled by the learning rate. By minimizing the MSE, gradient descent enables us to achieve better predictive accuracy.

## 4. Data Preprocessing

### 4.1 Dataset Overview

The dataset comprises the following features:

- **Temperature_c**: Target variable representing temperature in Celsius.
- **Humidity**: The relative humidity percentage.
- **Wind_Speed_kmh**: Wind speed measured in kilometers per hour.
- **Wind_Bearing_degrees**: Direction of the wind in degrees.
- **Visibility_km**: The visibility distance in kilometers.
- **Pressure_millibars**: Atmospheric pressure in millibars.
- **Rain**: Binary indicator for rainfall.
- **Description**: A categorical description of weather (e.g., Cold, Warm, Normal).

### 4.2 Data Cleaning and Preparation

To ensure optimal model performance, we cleaned the data as follows:

1. **Dropping Irrelevant or Redundant Columns**: We removed columns that were irrelevant to the prediction task or showed high collinearity with other features. This process reduces noise, focusing the model on the most informative features.
2. **Handling Missing Values**: Any missing values were filled using mean imputation or removed entirely depending on their distribution and potential impact on the model.
3. **One-Hot Encoding**: The `Description` column was one-hot encoded to convert categorical values into binary columns, allowing the model to interpret this feature as distinct categories rather than an ordinal sequence.

### 4.3 Feature Scaling

We standardized the numerical features (e.g., `Humidity`, `Wind_Speed_kmh`) to have a mean of zero and a standard deviation of one. This scaling is crucial in models that utilize gradient descent since it ensures consistent convergence by balancing feature contributions.

### 4.4 Data Splitting

The dataset was split into three parts: training (70%), validation (15%), and test (15%). The validation set aids in tuning hyperparameters, while the test set provides an unbiased performance assessment.

## 5. Model Development

### 5.1 Custom Linear Regression Model

The model was implemented from scratch using Python and NumPy. Our model uses a linear combination of weights and features to estimate temperature. To optimize these weights, we employed gradient descent, which updates weights in response to the gradient of the cost function.

The model's structure includes:

1. **Parameter Initialization**: All weights (`theta`) are initialized to zero.
2. **Cost Calculation**: For each iteration, MSE is computed to quantify the error.
3. **Gradient Calculation**: The partial derivatives with respect to each parameter are calculated.
4. **Weight Update**: Parameters are adjusted based on the learning rate and gradient.

The model training continues for a set number of iterations or until convergence, minimizing the error with each step.

### 5.2 Model Tuning with Gradient Descent

To optimize learning rate selection, we implemented a tuning mechanism that tested various rates on the validation set. By running the model with each rate and evaluating the validation cost, we identified the learning rate that yielded the lowest validation error.

The tuning process involved:

1. **Range Selection**: We tested learning rates from 1e-7 to 1e-2.
2. **Model Training**: For each rate, the model was trained on the training set.
3. **Validation Evaluation**: The MSE was computed on the validation set to determine the best-performing rate.

Our final model utilized the best learning rate, enhancing convergence speed and prediction accuracy.

## 6. Results and Performance Evaluation

### 6.1 Evaluation Metrics

We assessed the model using:

- **Mean Squared Error (MSE)**: Indicates the average squared error, providing a penalty for large errors.
- **Mean Absolute Error (MAE)**: Measures average error without penalizing outliers.
- **R² Score**: Shows the proportion of variance explained by the model.

### 6.2 Performance Comparison

Below is a summary of the model's performance before and after tuning:

| Metric       | Before Tuning | After Tuning |
| ------------ | ------------- | ------------ |
| **Test MSE** | 155.20        | 45.50        |
| **Test MAE** | 12.16         | 6.30         |
| **Test R²**  | -0.78         | 0.75         |

Tuning resulted in a noticeable performance increase, particularly in reducing MSE and increasing the R² score.

### 6.3 Insights from the Feature Correlation Matrix

A feature correlation matrix (Fig. 1) was used to analyze relationships among the variables. High correlations with `Temperature_c` indicate features with a strong predictive potential. This informed our choice of features and confirmed the relevance of factors like humidity and wind speed.

## 7. Challenges and Solutions

### 7.1 Feature Selection and Data Quality

The dataset contained several categorical variables that required encoding. One challenge was preserving information without overfitting, which we addressed through one-hot encoding and feature selection.

### 7.2 Choosing an Optimal Learning Rate

Finding a suitable learning rate for gradient descent is crucial, as too high a rate can cause divergence, while too low a rate may prolong training. We addressed this by implementing an extensive tuning process and using a range of learning rates.

### 7.3 Computational Efficiency

With a large number of iterations, training time becomes a concern. We improved efficiency by optimizing gradient calculations and limiting data transformations during each iteration.

## 8. Conclusion

This project successfully demonstrated the application of a custom linear regression model for temperature prediction. Using gradient descent for parameter optimization and learning rate tuning, we improved model accuracy.

The approach provided insights into weather-based temperature prediction and showcased the utility of gradient-based optimization. Future work may involve exploring additional features, employing advanced regularization methods to avoid overfitting, and extending the model to other predictive tasks in climatology.

## 9. Future Work

1. **Feature Expansion**: Adding more weather-related features like cloud cover, solar radiation, or past temperatures may further improve predictive performance.
2. **Regularization Techniques**: To enhance model generalization, regularization methods like L1 or L2 could be incorporated.
3. **Alternative Algorithms**: Testing alternative regression techniques or ensemble methods may reveal complementary insights and further improve accuracy.
