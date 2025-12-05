import {
  Platform,
  StyleSheet,
  View,
  Text,
  Pressable,
  KeyboardAvoidingView,
} from "react-native";
import { useState } from "react";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { HelloWave } from "@/components/hello-wave";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import { ProgressRing } from "@/components/ProgressRing";
import {
  getGreetingByName,
  postAgentsNutrition,
} from "@/lib/api/default/default";

export default function HomeScreen() {
  const [serverResponse, setServerResponse] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId] = useState(() => uuid7());

  const handleSubmitNutrition = async (text: string) => {
    setIsLoading(true);
    try {
      const response = await postAgentsNutrition({
        text,
        session_id: sessionId,
        user_id: "user-1",
      });
      console.log("Nutrition response:", response.data);
    } catch (error) {
      console.error("Error submitting nutrition:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const testServerConnection = async () => {
    try {
      const response = await getGreetingByName("John");
      if (response.status === 200) {
        setServerResponse(
          `Server says: ${JSON.stringify(response.data.message)}`,
        );
      }
    } catch (error) {
      setServerResponse(`Error: ${error}`);
      console.error("Error connecting to server:", error);
    }
  };

  return (
    <KeyboardAvoidingView behavior="height">
      <ThemedView className="pt-16 h-full w-full gap-4">
        <ThemedText className="px-4" type="title">
          Tuesday
        </ThemedText>
        <Pressable onPress={() => alert("Hello!")} className="px-4">
          <HelloWave />
        </Pressable>

        <Pressable
          onPress={testServerConnection}
          className="mx-4 p-4 bg-blue-500 rounded-lg"
        >
          <ThemedText className="text-white text-center">
            Test Server Connection
          </ThemedText>
        </Pressable>

        {serverResponse && (
          <ThemedView className="mx-4 p-4 bg-gray-100 rounded-lg">
            <ThemedText>{serverResponse}</ThemedText>
          </ThemedView>
        )}

        <TotalMacroPanel />
        <ThemedView className="flex-1 gap-4 p-4">
          <MealCard />
          <MealCard />
          <ProgressRing
            label="Calories"
            current={1200}
            goal={2000}
            color="#34D399"
            size="lg"
            unit="kcal"
          />
        </ThemedView>
        <InputBar
          onSubmitEditing={(event) =>
            handleSubmitNutrition(event.nativeEvent.text)
          }
        />
      </ThemedView>
    </KeyboardAvoidingView>
  );
}
