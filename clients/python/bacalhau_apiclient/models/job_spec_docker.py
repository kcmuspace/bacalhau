# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/filecoin-project/bacalhau.  # noqa: E501

    OpenAPI spec version: 0.3.18.post4
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


import pprint
import re  # noqa: F401

import six

from bacalhau_apiclient.configuration import Configuration


class JobSpecDocker(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    """
    Attributes:
      swagger_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    swagger_types = {
        'entrypoint': 'list[str]',
        'environment_variables': 'list[str]',
        'image': 'str',
        'working_directory': 'str'
    }

    attribute_map = {
        'entrypoint': 'Entrypoint',
        'environment_variables': 'EnvironmentVariables',
        'image': 'Image',
        'working_directory': 'WorkingDirectory'
    }

    def __init__(self, entrypoint=None, environment_variables=None, image=None, working_directory=None, _configuration=None):  # noqa: E501
        """JobSpecDocker - a model defined in Swagger"""  # noqa: E501
        if _configuration is None:
            _configuration = Configuration()
        self._configuration = _configuration

        self._entrypoint = None
        self._environment_variables = None
        self._image = None
        self._working_directory = None
        self.discriminator = None

        if entrypoint is not None:
            self.entrypoint = entrypoint
        if environment_variables is not None:
            self.environment_variables = environment_variables
        if image is not None:
            self.image = image
        if working_directory is not None:
            self.working_directory = working_directory

    @property
    def entrypoint(self):
        """Gets the entrypoint of this JobSpecDocker.  # noqa: E501

        optionally override the default entrypoint  # noqa: E501

        :return: The entrypoint of this JobSpecDocker.  # noqa: E501
        :rtype: list[str]
        """
        return self._entrypoint

    @entrypoint.setter
    def entrypoint(self, entrypoint):
        """Sets the entrypoint of this JobSpecDocker.

        optionally override the default entrypoint  # noqa: E501

        :param entrypoint: The entrypoint of this JobSpecDocker.  # noqa: E501
        :type: list[str]
        """

        self._entrypoint = entrypoint

    @property
    def environment_variables(self):
        """Gets the environment_variables of this JobSpecDocker.  # noqa: E501

        a map of env to run the container with  # noqa: E501

        :return: The environment_variables of this JobSpecDocker.  # noqa: E501
        :rtype: list[str]
        """
        return self._environment_variables

    @environment_variables.setter
    def environment_variables(self, environment_variables):
        """Sets the environment_variables of this JobSpecDocker.

        a map of env to run the container with  # noqa: E501

        :param environment_variables: The environment_variables of this JobSpecDocker.  # noqa: E501
        :type: list[str]
        """

        self._environment_variables = environment_variables

    @property
    def image(self):
        """Gets the image of this JobSpecDocker.  # noqa: E501

        this should be pullable by docker  # noqa: E501

        :return: The image of this JobSpecDocker.  # noqa: E501
        :rtype: str
        """
        return self._image

    @image.setter
    def image(self, image):
        """Sets the image of this JobSpecDocker.

        this should be pullable by docker  # noqa: E501

        :param image: The image of this JobSpecDocker.  # noqa: E501
        :type: str
        """

        self._image = image

    @property
    def working_directory(self):
        """Gets the working_directory of this JobSpecDocker.  # noqa: E501

        working directory inside the container  # noqa: E501

        :return: The working_directory of this JobSpecDocker.  # noqa: E501
        :rtype: str
        """
        return self._working_directory

    @working_directory.setter
    def working_directory(self, working_directory):
        """Sets the working_directory of this JobSpecDocker.

        working directory inside the container  # noqa: E501

        :param working_directory: The working_directory of this JobSpecDocker.  # noqa: E501
        :type: str
        """

        self._working_directory = working_directory

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value
        if issubclass(JobSpecDocker, dict):
            for key, value in self.items():
                result[key] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, JobSpecDocker):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, JobSpecDocker):
            return True

        return self.to_dict() != other.to_dict()
