import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import EditCollectionItem from '../form/EditCollectionItem';

const CollectionItemItem = ({data, index, editable, editing, checked, onToggleSelect, onChange, onDelete}) => {
  const item = data;

  if (editing) {
    return <EditCollectionItem data={data} index={index} onChange={onChange} onDelete={onDelete} />
  }

  const handleToggleSelect = () =>{
    onToggleSelect(index, data);
  }

  return (
    <tr className="">
      {editable && <td><input type="checkbox" onChange={handleToggleSelect} checked={checked} /></td>}
      <td>{ item.hash && <a href={`https://ipfs.io/ipfs/${item.hash}`} target="_blank"><p className="title">{item.hash}</p></a>}</td>
      <td><Link to={`/url?url=${item.url}`}>{item.url}</Link></td>
      <td>{item.description}</td>
    </tr>
  );
}

CollectionItemItem.propTypes = {

}

export default CollectionItemItem;
