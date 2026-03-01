import numpy as np
from scipy.spatial.distance import cdist


class CustomKMeans:
    def __init__(self, n_clusters: int = 2, max_iter: int = 100,
                 tol: float = 1e-4, random_state: int = 42,
                 init: str = 'k-means++'):
        self.n_clusters = n_clusters
        self.max_iter = max_iter
        self.tol = tol
        self.random_state = random_state
        self.init = init
        self.cluster_centers_ = None
        self.labels_ = None
        self.inertia_ = None
        self.n_iter_ = 0

        np.random.seed(random_state)

    def _initialize_centroids(self, X: np.ndarray) -> np.ndarray:
        n_samples, n_features = X.shape

        random_indices = np.random.choice(n_samples, self.n_clusters, replace=False)
        return X[random_indices]

    def _assign_clusters(self, X: np.ndarray, centroids: np.ndarray) -> np.ndarray:
        # calculate the distances from each point to each centroid
        distances = cdist(X, centroids)
        # Assigning a cluster placemark with a minimum distance
        return np.argmin(distances, axis=1)

    # Recalculating centroids as the midpoints of clusters.
    def _update_centroids(self, X: np.ndarray, labels: np.ndarray) -> np.ndarray:
        new_centroids = np.zeros((self.n_clusters, X.shape[1]))

        for i in range(self.n_clusters):
            # Get points of current class
            cluster_points = X[labels == i]

            if len(cluster_points) > 0:
                # The new centroid is the average of the cluster points
                new_centroids[i] = cluster_points.mean(axis=0)
            else:
                new_centroids[i] = X[np.random.randint(X.shape[0])]

        return new_centroids

    def fit(self, X: np.ndarray) -> 'CustomKMeans':
        X = np.asarray(X)
        n_samples, n_features = X.shape

        if n_samples < self.n_clusters:
            raise ValueError(f"n_samples={n_samples} should be >= n_clusters={self.n_clusters}")

        centroids = self._initialize_centroids(X)

        for iteration in range(self.max_iter):
            self.n_iter_ = iteration + 1
            labels = self._assign_clusters(X, centroids)
            new_centroids = self._update_centroids(X, labels)
            centroid_shift = np.sqrt(((new_centroids - centroids) ** 2).sum(axis=1)).max()
            centroids = new_centroids
            if centroid_shift < self.tol:
                print(f"Converged at iteration {iteration + 1}")
                break
        self.cluster_centers_ = centroids
        self.labels_ = self._assign_clusters(X, centroids)

        distances = cdist(X, self.cluster_centers_)
        self.inertia_ = np.sum(np.min(distances, axis=1) ** 2)

        return self

    def predict(self, X: np.ndarray) -> np.ndarray:
        if self.cluster_centers_ is None:
            raise ValueError("Model not fitted yet. Call fit() first.")

        return self._assign_clusters(X, self.cluster_centers_)

    def fit_predict(self, X: np.ndarray) -> np.ndarray:
        self.fit(X)
        return self.labels_