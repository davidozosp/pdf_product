from django.contrib import admin

# Register your models here.
from home.models import APIKey


class APIKeyAdmin(admin.ModelAdmin):
    list_display = ('key', 'requests', 'is_active')
    list_filter = ('is_active', 'date_joined', 'updated_at')
    search_fields = ('key',)
    
    
admin.site.register(APIKey, APIKeyAdmin)