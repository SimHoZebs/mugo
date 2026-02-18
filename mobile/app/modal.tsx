import { Link } from 'expo-router';
import { StyleSheet, View, Text } from 'react-native';

export default function ModalScreen() {
  return (
    <View style={styles.container} className="bg-stone-50 dark:bg-stone-950">
      <Text className="text-2xl font-bold leading-8 text-stone-950 dark:text-stone-50">This is a modal</Text>
      <Link href="/" dismissTo style={styles.link}>
        <Text className="text-base leading-[30px] text-cyan-700">Go to home screen</Text>
      </Link>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  link: {
    marginTop: 15,
    paddingVertical: 15,
  },
});
