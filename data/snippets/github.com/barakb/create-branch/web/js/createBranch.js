import { createBranchRequest, toggleXAPRequest } from "./actions";

export const CreateBranch = ({ name, onClick, onXAPToggle, isFetching, fromBranch, isXAPOnlyBranch }) => {
         let input;
         let from = fromBranch ? ("Create From " + fromBranch) : "Create From master"
         let btn = isFetching ? <button type="submit" className="btn btn-default disabled">{from}</button> :
            <button type="submit" className="btn btn-default" onClick={() => onClick(input.value, fromBranch, isXAPOnlyBranch)}>{from}</button>
         let checkbox = isXAPOnlyBranch ? <input type="checkbox" name="xap_only" onClick={onXAPToggle} checked>Create XAP Only Branch</input> :
             <input type="checkbox" name="xap_only" onClick={onXAPToggle}>Create XAP Only Branch</input>
         return (
             <form className="navbar-form navbar-left" role="search" onSubmit={ev => ev.preventDefault() } >
                 <div className="form-group">
                     <input type="text" className="form-control" placeholder="New Branch" ref={node => {input = node}}  defaultValue={name} />
                 </div>
                 {btn}
                  <div className="form-group">
                    {isFetching ? <img src="/web/images/gears.svg" className="form-control" /> : null}
                  </div>
                  <div>&nbsp;</div>
                 {checkbox}
             </form>
          )
}


const PropTypes = React.PropTypes;
CreateBranch.propTypes = {
  onClick: PropTypes.func.isRequired,
  onXAPToggle: PropTypes.func.isRequired,
  fromBranch : PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  isFetching: PropTypes.bool.isRequired,
  isXAPOnlyBranch: PropTypes.bool.isRequired,
};

const { Component } = React;

const mapStateToProps = (state) => {
    return {
        name : state.createBranch.name,
        fromBranch : state.createBranch.fromBranch,
        isFetching : state.createBranch.isFetching,
        isXAPOnlyBranch : state.createBranch.isXAPOnlyBranch,
    }
}

const mapDispatchToProps = (dispatch) => {
    return {
        onClick : (name, fromBranch, isXAPOnly) =>  { dispatch(createBranchRequest(name, fromBranch, isXAPOnly)); },
        onXAPToggle : () => { dispatch(toggleXAPRequest()); }
    }
}

const { connect } = ReactRedux;

export const CreateBranchContainer = connect(mapStateToProps, mapDispatchToProps)(CreateBranch);

