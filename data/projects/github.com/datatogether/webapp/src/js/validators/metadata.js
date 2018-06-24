import Validate from './index';

export default function validateMetadata(metadata = {}) {
  const errors = {};

  Object.keys(metadata.meta).forEach((fieldName) => {
    errors[`${fieldName}_field`] = Validate(fieldName).required().message();
    errors[`${fieldName}_value`] = Validate(metadata.meta[fieldName]).required().message();
  });

  errors.isValid = Object.keys(errors).every(key => (errors[key] == undefined));
  return errors;
}
