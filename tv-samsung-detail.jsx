<View>
  <Button text="hdmi" onPress={() => command('sendkey', { key: 'KEY_HDMI' })} />
  <Button text="tv" onPress={() => command('sendkey', { key: 'KEY_TV' })} />
</View>