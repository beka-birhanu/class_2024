import pandas as pd
import random
import os
import kaggle

kaggle.api.authenticate()

kaggle.api.dataset_download_files('nikhil7280/weather-type-classification', path=os.getcwd(), unzip=True)


# Read the input CSV file
df = pd.read_csv("weather_classification_data.csv")

# Define criteria for Good/Bad weather
def classify_weather(row):
    if (
        20 <= row["Temperature"] <= 30 and
        row["Humidity"] < 80 and
        row["Precipitation (%)"] < 20 and
        row["Cloud Cover"] in ["clear", "partly cloudy"]
    ):
        return "Good"
    else:
        return "Bad"

# Apply classification with a probabilistic override
def probabilistic_classify(row):
    true_classification = classify_weather(row)
    # Introduce a 10% chance of flipping the result
    if random.random() < 0.1:
        return "Good" if true_classification == "Bad" else "Bad"
    return true_classification

# Apply the probabilistic classification
df["Good/Bad"] = df.apply(probabilistic_classify, axis=1)

# Save the updated DataFrame to a new CSV file
df.to_csv("generated_weather_data.csv", index=False)
print("File saved to generated_weather_data.csv")

