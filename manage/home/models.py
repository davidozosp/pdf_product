from django.db import models

# Create your models here.


class APIKey(models.Model):
    key = models.CharField(max_length=100)
    requests = models.IntegerField(default=0)
    date_joined = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    is_active = models.BooleanField(default=True)
    def __str__(self):
        return self.key + " - " + str(self.requests)