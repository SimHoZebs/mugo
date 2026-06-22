import { by, device, expect, element } from 'detox';

describe('LazyFood App', () => {
  beforeAll(async () => {
    await device.launchApp();
  });

  beforeEach(async () => {
    await device.reloadReactNative();
  });

  it('should start and load the main layout', async () => {
    // This expects at least some basic structural element to be visible
    // Update the test IDs or text based on actual app content
    await expect(element(by.text('Mugo'))).toExist();
  });
});
