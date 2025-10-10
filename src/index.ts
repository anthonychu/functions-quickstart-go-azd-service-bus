import { app, InvocationContext } from '@azure/functions';

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

app.serviceBusQueue('serviceBusQueueTrigger', {
    queueName: '%ServiceBusQueueName%',
    connection: 'ServiceBusConnection',
    handler: async (message: unknown, context: InvocationContext): Promise<void> => {
        context.log('TypeScript ServiceBus Queue trigger start processing a message:', message);
        
        // Simulate the same 30-second processing time as the original Python function
        await delay(30000);
        
        context.log('TypeScript ServiceBus Queue trigger end processing a message');
    }
});