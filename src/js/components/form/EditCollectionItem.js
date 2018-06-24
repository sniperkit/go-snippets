import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import ValidInput from './ValidInput';
import ValidTextarea from './ValidTextarea';

const EditCollectionItem = ({data, index, onChange, onDelete}) => {
  const handleChange = (name, value) => {
    onChange(index, Object.assign({}, data, { [name] : value }));
  }

  const handleDelete = (e) => {
    onDelete(index, data);
  }

  return (
    <tr>
      <td><a onClick={handleDelete}>delete</a></td>
      <td></td>
      <td>{data.id ? <Link to={`/url?url=${item.url}`}>{item.url}</Link> : <ValidInput name="url" value={data.url} onChange={handleChange} />}</td>
      <td><ValidInput name="description" value={data.description} onChange={handleChange} /></td>
    </tr>
  );
}

EditCollectionItem.propTypes = {

}

export default EditCollectionItem;