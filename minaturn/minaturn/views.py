import json
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.shortcuts import get_object_or_404
from .models import Queue, QueueEntry


@csrf_exempt
def join_queue(request):
    if request.method == "POST":
        try:
            data = json.loads(request.body.decode("utf-8"))
            queue = get_object_or_404(Queue, id=data["queue_id"])

            entry = QueueEntry.objects.create(
                msisdn=data["msisdn"],
                queue=queue,
                left=False
            )
            return JsonResponse({"id": str(entry.id), "status": "joined"})
        except (KeyError, json.JSONDecodeError):
            return JsonResponse({"error": "Invalid request data"}, status=400)

    return JsonResponse({"error": "POST only"}, status=400)


def queue_position(request, queue_id, msisdn):
    try:
        entry = QueueEntry.objects.get(queue_id=queue_id, msisdn=msisdn, left=False)
    except QueueEntry.DoesNotExist:
        return JsonResponse({"error": "Not in queue"}, status=404)

    # Position = count of people who joined earlier and haven’t left
    position = (
        QueueEntry.objects.filter(queue_id=queue_id, left=False, joined_at__lt=entry.joined_at)
        .count()
        + 1
    )

    return JsonResponse({
        "msisdn": msisdn,
        "queue_id": str(queue_id),
        "position": position
    })


"""
For a given queue, we need the time to AVG_TIME_TO_PROCESS. 

Example:
 - A user joins the queue at 10:00 AM
 - The user is processed at 10:30 AM
 - The user gets served at 10:45 AM
Time to process for this user = 15 minutes

Do this for all users in the queue IN THE LAST 30 MINUTES
This will help us understand the average time to process a user in the queue.
This can be used to estimate the time for the next user in the queue.


"""