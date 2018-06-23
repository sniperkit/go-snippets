# Create your views here.

from django.http import HttpResponse
from django.views.generic import TemplateView
from random import randint

class DisplayAdditionView(TemplateView):
	#`template_name = "addition.html"
	def get_context_data(self, **kwargs):
		#return HttpResponse(str(randint(0, 10))+'+'+str(randint(0,10)))
		context = super(DispalyAdditionView, self).get_context_data(**kwargs)
		#return TemplateView(request, template="addition.html")
		return context
