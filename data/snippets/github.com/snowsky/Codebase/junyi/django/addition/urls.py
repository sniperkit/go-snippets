from django.conf.urls import patterns, include, url
from django.views.generic import TemplateView
from addition.views import DisplayAdditionView

# Uncomment the next two lines to enable the admin:
urlpatterns = patterns('',
    #url(r'^', ListView),
    url(r'^$', TemplateView.as_view(template_name='addition.html')),
)
