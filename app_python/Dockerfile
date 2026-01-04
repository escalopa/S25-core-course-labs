# Use specific Alpine-based Python image for smaller size
FROM python:3.11-alpine3.19

# Set working directory
WORKDIR /app

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

# Copy requirements first for better layer caching
COPY --chown=appuser:appuser requirements.txt .

# Install dependencies as root (needed for system packages)
RUN pip install --no-cache-dir -r requirements.txt

# Switch to non-root user
USER appuser

# Copy only necessary application files
COPY --chown=appuser:appuser app.py .
COPY --chown=appuser:appuser templates ./templates

# Expose port
EXPOSE 5000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5000/health || exit 1

# Run the application
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "5000"]
