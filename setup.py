import pathlib
from setuptools import setup, find_packages

# The directory containing this file
HERE = pathlib.Path(__file__).parent

# The text of the README file
README = (HERE / "README.md").read_text()

setup(
     name='dnsexit-ip-update',  
     version='0.1.1',
     author="Christian Tragesser",
     author_email="christian@evoen.net",
     description="Dynamic DNS client for DNS Exit",
     long_description_content_type="text/markdown",
     long_description=README,
     license='MIT',
     url="https://github.com/christianTragesser/dnsexit-ip-update",
     packages=find_packages(exclude=["tests"]),
     install_requires=[
        "requests >= 2.24.0",
        "python-json-logger >= 0.1.11",
        "dnspython >= 1.16.0"
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
 )
