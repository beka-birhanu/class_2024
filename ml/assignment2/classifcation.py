import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder, StandardScaler
from sklearn.linear_model import LogisticRegression
from sklearn.tree import DecisionTreeClassifier
from sklearn.ensemble import RandomForestClassifier
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score, classification_report, confusion_matrix

# Load the dataset
data = pd.read_csv("generated_weather_data.csv")  # Replace with your actual file path

# Preview the data
print("Raw Data: \n", data.head())

# Preprocessing
# Encode categorical variables
label_encoders = {}
for col in ['Cloud Cover', 'Season', 'Location', 'Weather Type', 'Good/Bad']:
    label_encoders[col] = LabelEncoder()
    data[col] = label_encoders[col].fit_transform(data[col])

print("\nProcessed Data:\n", data.head())

# Separate features and target variable
X = data.drop(columns=['Weather Type'])  # Features
y = data['Weather Type']                # Target

# Scale only numerical features
numerical_features = ['Temperature', 'Humidity', 'Wind Speed', 'Precipitation (%)', 
                      'Atmospheric Pressure', 'UV Index', 'Visibility (km)']
X[numerical_features] = StandardScaler().fit_transform(X[numerical_features])

# Split the data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Train and Evaluate Models
models = {
    "Logistic Regression": LogisticRegression(max_iter=1000),
    "Decision Tree": DecisionTreeClassifier(random_state=42),
    "Random Forest": RandomForestClassifier(random_state=42, n_estimators=100),
    "SVM": SVC(kernel='linear', random_state=42)
}

results = {}

for name, model in models.items():
    # Train the model
    model.fit(X_train, y_train)
    
    # Make predictions
    y_pred = model.predict(X_test)
    
    # Evaluate the model
    accuracy = accuracy_score(y_test, y_pred)
    precision = precision_score(y_test, y_pred, average="weighted")
    recall = recall_score(y_test, y_pred, average="weighted")
    f1 = f1_score(y_test, y_pred, average="weighted")
    conf_matrix = confusion_matrix(y_test, y_pred)
    class_report = classification_report(y_test, y_pred)
    
    # Store results
    results[name] = {
        "Accuracy": accuracy,
        "Precision": precision,
        "Recall": recall,
        "F1-Score": f1,
        "Confusion Matrix": conf_matrix,
        "Classification Report": class_report
    }


# Print detailed results
print("\nDetailed Results:")
for name, metrics in results.items():
    print(f"\nModel: {name}")
    print(f"Accuracy: {metrics['Accuracy']:.2f}")
    print(f"Precision: {metrics['Precision']:.2f}")
    print(f"Recall: {metrics['Recall']:.2f}")
    print(f"F1-Score: {metrics['F1-Score']:.2f}")
    print("\nConfusion Matrix:")
    print(metrics["Confusion Matrix"])
    print("\nClassification Report:")
    print(metrics["Classification Report"])

# Create a summary table for key metrics
summary = {
    "Model": [],
    "Accuracy": [],
    "Precision": [],
    "Recall": [],
    "F1-Score": []
}

for name, metrics in results.items():
    summary["Model"].append(name)
    summary["Accuracy"].append(metrics["Accuracy"])
    summary["Precision"].append(metrics["Precision"])
    summary["Recall"].append(metrics["Recall"])
    summary["F1-Score"].append(metrics["F1-Score"])

summary_df = pd.DataFrame(summary)

# Print the summary table
print("\nModel Performance Summary:")
print(summary_df)

# Plot the summary
summary_df.set_index("Model", inplace=True)

# Plot each metric
summary_df.plot(kind="bar", figsize=(10, 6), colormap="viridis", edgecolor="black")
plt.title("Model Performance Summary", fontsize=16)
plt.ylabel("Score", fontsize=14)
plt.xlabel("Model", fontsize=14)
plt.xticks(rotation=45, fontsize=12)
plt.legend(loc="best", fontsize=12, title="Metrics")

# Set y-axis ticks 0.05 apart
y_ticks = np.arange(0, 1.05, 0.05)  # From 0 to 1.0 in steps of 0.05
plt.yticks(y_ticks, fontsize=12)
plt.grid(axis="y", linestyle="--", alpha=0.7)

# Save or show the plot
plt.tight_layout()
plt.show()
