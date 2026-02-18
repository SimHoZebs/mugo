import useGlobalStore from "@/lib/store";

describe("useGlobalStore", () => {
  beforeEach(() => {
    useGlobalStore.setState({
      meals: [],
      userProfile: null,
    });
  });

  describe("authentication", () => {
    it("should login and set userProfile", () => {
      const { login } = useGlobalStore.getState();

      login("testuser");

      expect(useGlobalStore.getState().userProfile).toEqual({
        username: "testuser",
        name: "",
        dietaryPreferences: [],
        weight: 0,
        height: 0,
        unitSystem: "metric",
      });
    });

    it("should logout and clear userProfile", () => {
      const { login, logout } = useGlobalStore.getState();

      login("testuser");
      expect(useGlobalStore.getState().userProfile).not.toBeNull();

      logout();
      expect(useGlobalStore.getState().userProfile).toBeNull();
    });
  });

  describe("meals", () => {
    const mockMeal = {
      id: "meal-1",
      createdAt: "2024-01-01T00:00:00Z",
      nutrition: {
        name: "Test Meal",
        macros: {
          calories: 500,
          protein: 30,
          carbs: 50,
          fat: 20,
        },
        assumptions: [],
      },
    };

    it("should set meals", () => {
      const { setMeals } = useGlobalStore.getState();

      setMeals([mockMeal as any]);

      expect(useGlobalStore.getState().meals).toHaveLength(1);
      expect(useGlobalStore.getState().meals[0].id).toBe("meal-1");
    });

    it("should update a specific meal", () => {
      const { setMeals, updateMeal } = useGlobalStore.getState();

      setMeals([mockMeal as any]);
      updateMeal("meal-1", {
        ...mockMeal,
        nutrition: { ...mockMeal.nutrition, name: "Updated Meal" },
      } as any);

      expect(useGlobalStore.getState().meals[0].nutrition.name).toBe(
        "Updated Meal"
      );
    });
  });

  describe("userProfile updates", () => {
    it("should update userProfile fields", () => {
      const { login, updateUserProfile } = useGlobalStore.getState();

      login("testuser");
      updateUserProfile({ name: "Test User", weight: 70 });

      const profile = useGlobalStore.getState().userProfile;
      expect(profile?.name).toBe("Test User");
      expect(profile?.weight).toBe(70);
    });
  });

  describe("dietary preferences", () => {
    it("should add dietary preference", () => {
      const { login, addDietaryPreference } = useGlobalStore.getState();

      login("testuser");
      addDietaryPreference({ id: "pref-1", text: "Vegetarian" });

      const profile = useGlobalStore.getState().userProfile;
      expect(profile?.dietaryPreferences).toHaveLength(1);
      expect(profile?.dietaryPreferences[0].text).toBe("Vegetarian");
    });

    it("should update dietary preference", () => {
      const { login, addDietaryPreference, updateDietaryPreference } =
        useGlobalStore.getState();

      login("testuser");
      addDietaryPreference({ id: "pref-1", text: "Vegetarian" });
      updateDietaryPreference("pref-1", "Vegan");

      const profile = useGlobalStore.getState().userProfile;
      expect(profile?.dietaryPreferences[0].text).toBe("Vegan");
    });

    it("should remove dietary preference", () => {
      const { login, addDietaryPreference, removeDietaryPreference } =
        useGlobalStore.getState();

      login("testuser");
      addDietaryPreference({ id: "pref-1", text: "Vegetarian" });
      addDietaryPreference({ id: "pref-2", text: "Gluten-free" });
      removeDietaryPreference("pref-1");

      const profile = useGlobalStore.getState().userProfile;
      expect(profile?.dietaryPreferences).toHaveLength(1);
      expect(profile?.dietaryPreferences[0].text).toBe("Gluten-free");
    });
  });
});
