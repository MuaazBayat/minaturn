from django.contrib import admin
from .models import Queue, QueueEntry


@admin.register(Queue)
class QueueAdmin(admin.ModelAdmin):
    list_display = ("id", "name", "description", "created_at")
    search_fields = ("name",)
    ordering = ("created_at",)


@admin.register(QueueEntry)
class QueueEntryAdmin(admin.ModelAdmin):
    list_display = ("msisdn", "queue", "joined_at", "left")
    list_filter = ("left", "queue")
    search_fields = ("msisdn",)
    ordering = ("joined_at",)
