import sys
import json
import numpy as np
import pandas as pd
from sklearn.cluster import KMeans

def main():
    if len(sys.argv) != 2:
        print("Usage: python algorithm.py <dataset_file>")
        sys.exit(1)

    dataset_file = sys.argv[1]

    try:
        # Load dataset
        data = pd.read_csv(dataset_file)

        # Ensure dataset is numeric
        if not np.issubdtype(data.dtypes.values[0], np.number):
            raise ValueError("Dataset must contain only numeric values.")

        # Convert to NumPy array
        X = data.values

        # Apply k-means clustering
        n_clusters = 3  # You can adjust the number of clusters
        kmeans = KMeans(n_clusters=n_clusters, random_state=42)
        kmeans.fit(X)

        # Extract results
        centroids = np.round(kmeans.cluster_centers_, 4).tolist()
        inertia = round(kmeans.inertia_, 4)

        # Create output JSON
        output = {
            "result": {
                "centroids": centroids,
                "inertia": inertia
            }
        }

        # Print output as JSON
        print(json.dumps(output))

    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()