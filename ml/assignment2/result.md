#### **Data Overview**
**Raw Data (Sample):**
| Temperature | Humidity | Wind Speed | Precipitation (%) | Cloud Cover    | Season  | Visibility (km) | Location  | Weather Type | Good/Bad |
|-------------|----------|------------|-------------------|----------------|---------|-----------------|-----------|--------------|----------|
| 14.0        | 73       | 9.5        | 82.0              | partly cloudy | Winter  | 3.5             | inland    | Rainy        | Bad      |
| 39.0        | 96       | 8.5        | 71.0              | partly cloudy | Spring  | 10.0            | inland    | Cloudy       | Bad      |
| 30.0        | 64       | 7.0        | 16.0              | clear          | Spring  | 5.5             | mountain  | Sunny        | Good     |
| 38.0        | 83       | 1.5        | 82.0              | clear          | Spring  | 1.0             | coastal   | Sunny        | Bad      |
| 27.0        | 74       | 17.0       | 66.0              | overcast       | Winter  | 2.5             | mountain  | Rainy        | Bad      |

---

**Processed Data (Sample):**
| Temperature | Humidity | Wind Speed | Precipitation (%) | Cloud Cover   | Season | Visibility (km) | Location | Weather Type | Good/Bad |
|-------------|----------|------------|-------------------|-------------|--------|-----------------|----------|--------------|----------|
| 14.0        | 73       | 9.5        | 82.0              | 3            | 3      | 3.5             | 1        | 1            | 0        |
| 39.0        | 96       | 8.5        | 71.0              | 3            | 1      | 10.0            | 1        | 0            | 0        |
| 30.0        | 64       | 7.0        | 16.0              | 0            | 1      | 5.5             | 2        | 3            | 1        |
| 38.0        | 83       | 1.5        | 82.0              | 0            | 1      | 1.0             | 0        | 3            | 0        |
| 27.0        | 74       | 17.0       | 66.0              | 2            | 3      | 2.5             | 2        | 1            | 0        |

---

### **Model Performances**

#### Logistic Regression
- **Accuracy**: 85.26%
- **Precision**: 85.24%
- **Recall**: 85.27%
- **F1-Score**: 85.21%

**Confusion Matrix**:
| Predicted\Actual | 0   | 1   | 2   | 3   |
|-------------------|-----|-----|-----|-----|
| **0**            | 532 | 61  | 18  | 40  |
| **1**            | 35  | 541 | 46  | 25  |
| **2**            | 21  | 6   | 658 | 16  |
| **3**            | 67  | 38  | 16  | 520 |

---

#### Decision Tree
- **Accuracy**: 90.72%
- **Precision**: 90.73%
- **Recall**: 90.72%
- **F1-Score**: 90.71%

**Confusion Matrix**:
| Predicted\Actual | 0   | 1   | 2   | 3   |
|-------------------|-----|-----|-----|-----|
| **0**            | 581 | 33  | 20  | 17  |
| **1**            | 36  | 575 | 15  | 21  |
| **2**            | 18  | 10  | 663 | 10  |
| **3**            | 24  | 18  | 23  | 576 |

---

#### Random Forest
- **Accuracy**: 91.40%
- **Precision**: 91.46%
- **Recall**: 91.40%
- **F1-Score**: 91.41%

**Confusion Matrix**:
| Predicted\Actual | 0   | 1   | 2   | 3   |
|-------------------|-----|-----|-----|-----|
| **0**            | 591 | 28  | 19  | 13  |
| **1**            | 33  | 589 | 12  | 13  |
| **2**            | 21  | 9   | 659 | 12  |
| **3**            | 30  | 16  | 21  | 574 |

---

#### SVM
- **Accuracy**: 87.95%
- **Precision**: 87.92%
- **Recall**: 87.95%
- **F1-Score**: 87.93%

**Confusion Matrix**:
| Predicted\Actual | 0   | 1   | 2   | 3   |
|-------------------|-----|-----|-----|-----|
| **0**            | 549 | 51  | 13  | 38  |
| **1**            | 27  | 560 | 34  | 26  |
| **2**            | 13  | 3   | 661 | 24  |
| **3**            | 52  | 26  | 11  | 552 |

---

### **Performance Summary Table**
| Model                | Accuracy | Precision | Recall | F1-Score |
|-----------------------|----------|-----------|--------|----------|
| Logistic Regression   | 0.852652 | 0.852352  | 0.852652 | 0.852128 |
| Decision Tree         | 0.907197 | 0.907270  | 0.907197 | 0.907120 |
| Random Forest         | 0.914015 | 0.914585  | 0.914015 | 0.914093 |
| SVM                   | 0.879545 | 0.879168  | 0.879545 | 0.879298 |
