import pandas as pd
import random

# Read the input CSV file
input_file = "weather_data.csv"  
df = pd.read_csv(input_file)

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
output_file = "generated_weather_data.csv"  # Path to save the new CSV file
df.to_csv(output_file, index=False)

print(f"File saved to {output_file}")

