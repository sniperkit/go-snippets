import React, { PropTypes } from 'react';

import List from '../List';
import ValidInput from './ValidInput';
import SelectUser from './SelectUser';
import ValidTextarea from './ValidTextarea';

const CollectionForm = ({ name, data, users, onChange, onCancel, onSubmit }) => {
  const collection = data;
  const handleChange = (name, value) => {
    onChange(name, Object.assign({}, data, { [name] : value }));
  }
  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(collection, e);
  };

  const handleItemChange = (index, item) => {
    let contents = collection.contents.concat()
    contents.splice(index,1,item);

    onChange(name, Object.assign({}, collection, { contents }));
  }

  const handleItemDelete = (index) => {
    // onChange(name, data);
  }

  const handleAddItem = (e) => {
    e.preventDefault();
    let contents = collection.contents || [];
    onChange(name,Object.assign({}, collection, { 
      contents: contents.concat([
        { url : "", description: "", index: -1 }
      ]),
    }));
  }

  return (
    <div id="collection" className="page">
      <div className="container">
        <header className="row">
          <div className="col-md-12">
            <hr className="green" />
            <label className="label">Collection</label>
          </div>
        </header>
        <div className="row">
          <div className="col-md-12">
            <SelectUser label="owner" name="creator" users={users} value={collection.creator} onChange={handleChange} />
            <ValidInput name="title" label="title" value={collection.title} onChange={handleChange} />
            <ValidTextarea name="description" label="description" value={collection.description} onChange={handleChange} />
            <br />
            <button className="btn" onClick={onCancel}>Cancel</button>
            <input className="btn btn-primary" type="submit" value="Save" onClick={handleSubmit} />
          </div>
        </div>
      </div>
    </div>
  );
};

CollectionForm.propTypes = {
  data: PropTypes.object.isRequired,
  name: PropTypes.string,
  onChange: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

CollectionForm.defaultProps = {
  name : "collection",
}

export default CollectionForm;
