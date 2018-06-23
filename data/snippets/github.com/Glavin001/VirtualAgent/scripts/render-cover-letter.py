# http://eosrei.net/articles/2015/11/latex-templates-python-and-jinja2-generate-pdfs

import jinja2
import os
from jinja2 import Template
from subprocess import call

tex_template_path = os.path.abspath('template/cover-letter')
print(tex_template_path)
latex_jinja_env = jinja2.Environment(
    block_start_string='\BLOCK{',
    block_end_string='}',
    variable_start_string='\VAR{',
    variable_end_string='}',
    comment_start_string='\#{',
    comment_end_string='}',
    line_statement_prefix='%=',
    line_comment_prefix='%#',
    trim_blocks=True,
    autoescape=False,
    loader=jinja2.FileSystemLoader(tex_template_path)
)

username = 'glavin001'

# Make User's Directory
user_dir = os.path.abspath('build/' + username)
if not os.path.exists(user_dir):
    os.makedirs(user_dir)
    # Create symbolic links
    # os.symlink(os.path.abspath('template/awesome-cv.cls'), os.path.abspath('build/'+username+'/awesome-cv.cls'))

userData = {
  "basics": {
    "name": "John Doe",
    "label": "Programmer",
    "picture": "",
    "email": "john@gmail.com",
    "phone": "(912) 555-4321",
    "website": "http://johndoe.com",
    "summary": "A summary of John Doe...",
    "location": {
      "address": "2712 Broadway St",
      "postalCode": "CA 94115",
      "city": "San Francisco",
      "countryCode": "US",
      "region": "California"
    },
    "profiles": [{
      "network": "Twitter",
      "username": "john",
      "url": "http://twitter.com/john"
    }]
  },
  "work": [{
    "company": "Company",
    "position": "President",
    "website": "http://company.com",
    "startDate": "2013-01-01",
    "endDate": "2014-01-01",
    "summary": "Description...",
    "highlights": [
      "Started the company"
    ]
  }],
  "volunteer": [{
    "organization": "Organization",
    "position": "Volunteer",
    "website": "http://organization.com/",
    "startDate": "2012-01-01",
    "endDate": "2013-01-01",
    "summary": "Description...",
    "highlights": [
      "Awarded 'Volunteer of the Month'"
    ]
  }],
  "education": [{
    "institution": "University",
    "area": "Software Development",
    "studyType": "Bachelor",
    "startDate": "2011-01-01",
    "endDate": "2013-01-01",
    "gpa": "4.0",
    "courses": [
      "DB1101 - Basic SQL"
    ]
  }],
  "awards": [{
    "title": "Award",
    "date": "2014-11-01",
    "awarder": "Company",
    "summary": "There is no spoon."
  }],
  "publications": [{
    "name": "Publication",
    "publisher": "Company",
    "releaseDate": "2014-10-01",
    "website": "http://publication.com",
    "summary": "Description..."
  }],
  "skills": [{
    "name": "Web Development",
    "level": "Master",
    "keywords": [
      "HTML",
      "CSS",
      "Javascript"
    ]
  }],
  "languages": [{
    "name": "English",
    "level": "Native speaker"
  }],
  "interests": [{
    "name": "Wildlife",
    "keywords": [
      "Ferrets",
      "Unicorns"
    ]
  }],
  "references": [{
    "name": "Jane Doe",
    "reference": "Reference..."
  }]
}
jobData = {
    'title': 'Front End Web Developer',
    'source': 'GitHub Careers',
    'company': {
        'name': 'GitHub'
    }
}

# Render 1 - Intro
template = latex_jinja_env.get_template('1-intro/1.tex')
intro = template.render(user=userData, job=jobData)

# Render 2 - Company Fit
template = latex_jinja_env.get_template('2-company_fit/1.tex')
companyFit = template.render(user=userData, job=jobData)

# Render 3 - Skills & Projects
template = latex_jinja_env.get_template('3-skills_projects/1.tex')
skillsProjects = template.render(user=userData, job=jobData)

# Render 4 - Work Fit
template = latex_jinja_env.get_template('4-work_fit/1.tex')
workFit = template.render(user=userData, job=jobData)

# Render 5 - Volunteer
template = latex_jinja_env.get_template('5-volunteer/1.tex')
volunteer = template.render(user=userData, job=jobData)

# Render 6 - Thanks
template = latex_jinja_env.get_template('6-thanks/1.tex')
thanks = template.render(user=userData, job=jobData)

# Render 7 - Complete Cover Letter
template = latex_jinja_env.get_template('template.tex')
cover_letter = template.render(
    user=userData,
    job=jobData,
    paragraphs={
        'intro': intro,
        'company': companyFit,
        'skills': skillsProjects,
        'work': workFit,
        'volunteer': volunteer,
        'thanks': thanks
    }
)
# print(cover_letter)

# Write to file
f = open(os.path.abspath('build/'+username+'/cover-letter.tex'),'w')
f.write(cover_letter.encode('utf-8')) # python will convert \n to os.linesep
f.close() # you can omit in most cases as the destructor will call it

# Build PDF
call(["xelatex","-output-directory="+os.path.abspath('build/'+username+'/'),"cover-letter.tex"], cwd=os.path.abspath('template/'))
call(["convert","-density", "300", os.path.abspath('build/'+username+'/cover-letter.pdf'), "-quality", "100", os.path.abspath('build/'+username+'/cover-letter.png')])
