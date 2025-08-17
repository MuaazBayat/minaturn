from django.db import models
import uuid

class Queue(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    name = models.CharField(max_length=255)
    description = models.TextField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return self.name


class QueueEntry(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    msisdn = models.CharField(max_length=15)  # phone number
    joined_at = models.DateTimeField(auto_now_add=True)
    left = models.BooleanField(default=False)
    queue = models.ForeignKey(Queue, on_delete=models.CASCADE, related_name='items')

    def __str__(self):
        return f"{self.msisdn} in {self.queue.name}"
