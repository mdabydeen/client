// @flow
import React, {Component} from 'react'
import type {Props} from './index.render'
import {Box, Checkbox, Button, Text, Icon} from '../../../common-adapters'
import {globalColors, globalStyles, globalMargins} from '../../../styles'
import {getStyle} from '../../../common-adapters/text'

/* types:
  paperkey: HiddenString,
  onFinish: () => void,
  onBack: () => void,
  title?: ?string
  */

type State = {
  checked: boolean,
}

class SuccessRender extends Component<void, Props, State> {
  state: State;

  constructor (props: Props) {
    super(props)
    this.state = {
      checked: false,
    }
  }

  render () {
    return (
      <Box style={{padding: globalMargins.large, flex: 1}}>
        <Text type='Header' style={textCenter}>{this.props.title || "Congratulations, you've just joined Keybase!"}</Text>
        <Text type='Body' style={{...textCenter, marginTop: globalMargins.medium}}>Here is your unique paper key, it will allow you to perform important Keybase tasks in the future. This is the only time you'll see this so be sure to write it down.</Text>

        <Box style={paperKeyContainerStyle}>
          <Text type='Header' style={paperkeyStyle}>{this.props.paperkey.stringValue()}</Text>
          <Icon type='icon-paper-key-corner' style={paperCornerStyle} />
        </Box>

        <Box style={confirmCheckboxStyle}>
          <Checkbox
            label='Yes, I wrote this down.'
            checked={this.state.checked}
            onCheck={checked => this.setState({checked})} />
        </Box>

        <Box style={{flex: 2, justifyContent: 'flex-end'}}>
          <Button style={buttonStyle}
            disabled={!this.state.checked}
            waiting={this.props.waiting}
            onClick={this.props.onFinish}
            label='Done'
            type='Primary' />
        </Box>
      </Box>
    )
  }
}

const confirmCheckboxStyle = {
  ...globalStyles.flexBoxRow,
  alignSelf: 'center',
}

const buttonStyle = {
}

const textCenter = {
  textAlign: 'center',
}

const paperKeyContainerStyle = {
  position: 'relative',
  alignSelf: 'center',
  marginTop: globalMargins.large,
  marginBottom: globalMargins.large,
  paddingTop: globalMargins.small,
  paddingBottom: globalMargins.small,
  paddingLeft: globalMargins.small,
  paddingRight: globalMargins.large,
  borderRadius: 4,
  backgroundColor: globalColors.white,
  borderStyle: 'solid',
  borderWidth: 4,
  borderColor: globalColors.darkBlue,
}

const paperkeyStyle = {
  ...getStyle('Header', 'Normal'),
  ...globalStyles.fontTerminal,
  color: globalColors.darkBlue,
  textAlign: 'center',
}

const paperCornerStyle = {
  position: 'absolute',
  right: 0,
  top: -4,
}

export default SuccessRender
